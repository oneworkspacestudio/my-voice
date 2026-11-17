// Resident CPU Qwen adapter for the exact llama.cpp revision in runtime.lock.json.
// Public C APIs only; symbols are resolved from the bundled runtime, never PATH.
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include "mtmd-helper.h"
#include <stdexcept>
#include <cstring>

#define FUNCTIONS(X) \
 X(ggml_backend_load_all_from_path) \
 X(llama_backend_init) X(llama_model_default_params) X(llama_context_default_params) \
 X(llama_model_load_from_file) X(llama_init_from_model) X(llama_model_free) X(llama_free) \
 X(llama_get_memory) X(llama_memory_clear) X(llama_get_embeddings_ith) X(llama_log_set) \
 X(llama_sampler_chain_default_params) X(llama_sampler_chain_init) X(llama_sampler_chain_add) \
 X(llama_sampler_init_top_k) X(llama_sampler_init_top_p) X(llama_sampler_init_temp) \
 X(llama_sampler_init_dist) X(llama_sampler_sample) X(llama_sampler_reset) X(llama_sampler_free) \
 X(mtmd_context_params_default) X(mtmd_init_from_file) X(mtmd_free) X(mtmd_bitmap_free) \
 X(mtmd_helper_init_opt_default) X(mtmd_helper_bitmap_init_from_file) X(mtmd_helper_log_set) \
 X(mtmd_helper_gen_audio_init) X(mtmd_helper_gen_audio_free) X(mtmd_helper_gen_audio_reset) \
 X(mtmd_helper_gen_audio_inp_default) X(mtmd_helper_gen_audio_set_input) \
 X(mtmd_helper_gen_audio_step_prompt) X(mtmd_helper_gen_audio_step_gen) \
 X(mtmd_helper_gen_audio_get_output) X(mtmd_helper_gen_audio_flush)
#define DECLARE(n) decltype(&n) p_##n;
FUNCTIONS(DECLARE)
static void quiet(ggml_log_level, const char *, void *) {}
static const char *last_error = "initialization";
extern "C" __declspec(dllexport) const char *voice_last_error() { return last_error; }

struct Voice {
    llama_model *model = nullptr;
    llama_context *ctx = nullptr;
    llama_sampler *sampler = nullptr;
    mtmd_context *codec = nullptr;
    mtmd_bitmap *speaker = nullptr;
    mtmd_helper_gen_audio *gen = nullptr;
    ~Voice() {
        if (gen) p_mtmd_helper_gen_audio_free(gen);
        if (speaker) p_mtmd_bitmap_free(speaker);
        if (codec) p_mtmd_free(codec);
        if (sampler) p_llama_sampler_free(sampler);
        if (ctx) p_llama_free(ctx);
        if (model) p_llama_model_free(model);
    }
};

extern "C" __declspec(dllexport) void *voice_create(const wchar_t *dir, const char *model,
        const char *codec, const char *reference, int threads) {
    Voice *v = nullptr;
    try {
        static bool initialized = false;
        if (!initialized) {
            SetDefaultDllDirectories(LOAD_LIBRARY_SEARCH_DEFAULT_DIRS);
            AddDllDirectory(dir);
            const auto llama = LoadLibraryExW((std::wstring(dir)+L"/llama.dll").c_str(), nullptr, LOAD_LIBRARY_SEARCH_DEFAULT_DIRS);
            const auto mtmd = LoadLibraryExW((std::wstring(dir)+L"/mtmd.dll").c_str(), nullptr, LOAD_LIBRARY_SEARCH_DEFAULT_DIRS);
            const auto ggml = LoadLibraryExW((std::wstring(dir)+L"/ggml.dll").c_str(), nullptr, LOAD_LIBRARY_SEARCH_DEFAULT_DIRS);
            if (!llama || !mtmd || !ggml) { last_error = "bundled library load"; return nullptr; }
#define RESOLVE(n) p_##n = reinterpret_cast<decltype(p_##n)>(GetProcAddress(llama, #n)); \
            if (!p_##n) p_##n = reinterpret_cast<decltype(p_##n)>(GetProcAddress(mtmd, #n)); \
            if (!p_##n) p_##n = reinterpret_cast<decltype(p_##n)>(GetProcAddress(ggml, #n)); \
            if (!p_##n) { last_error = #n; throw std::runtime_error("Bundled ABI mismatch"); }
            FUNCTIONS(RESOLVE)
            p_llama_log_set(quiet, nullptr);
            p_mtmd_helper_log_set(quiet, nullptr);
            auto length = WideCharToMultiByte(CP_UTF8,0,dir,-1,nullptr,0,nullptr,nullptr);
            std::string backend_path(length,'\0');
            WideCharToMultiByte(CP_UTF8,0,dir,-1,backend_path.data(),length,nullptr,nullptr);
            p_ggml_backend_load_all_from_path(backend_path.c_str());
            p_llama_backend_init();
            initialized = true;
        }
        v = new Voice;
        auto mp = p_llama_model_default_params();
        mp.n_gpu_layers = 0;
        v->model = p_llama_model_load_from_file(model, mp);
        if (!v->model) { last_error = "model"; throw std::runtime_error("model"); }
        auto cp = p_llama_context_default_params();
        cp.n_ctx = 4096; cp.n_batch = 512; cp.n_ubatch = 512;
        cp.n_threads = cp.n_threads_batch = threads;
        cp.embeddings = true; cp.offload_kqv = false;
        v->ctx = p_llama_init_from_model(v->model, cp);
        if (!v->ctx) { last_error = "context"; throw std::runtime_error("context"); }
        auto tp = p_mtmd_context_params_default();
        tp.use_gpu = false; tp.n_threads = threads;
        v->codec = p_mtmd_init_from_file(codec, v->model, tp);
        if (!v->codec) { last_error = "codec"; throw std::runtime_error("codec"); }
        v->speaker = p_mtmd_helper_bitmap_init_from_file(v->codec, reference, false, p_mtmd_helper_init_opt_default()).bitmap;
        v->gen = p_mtmd_helper_gen_audio_init(v->ctx, v->codec);
        if (!v->speaker || !v->gen) { last_error = "speaker"; throw std::runtime_error("speaker"); }
        v->sampler = p_llama_sampler_chain_init(p_llama_sampler_chain_default_params());
        p_llama_sampler_chain_add(v->sampler, p_llama_sampler_init_top_k(40));
        p_llama_sampler_chain_add(v->sampler, p_llama_sampler_init_top_p(0.9f, 1));
        p_llama_sampler_chain_add(v->sampler, p_llama_sampler_init_temp(0.8f));
        p_llama_sampler_chain_add(v->sampler, p_llama_sampler_init_dist(42));
        return v;
    } catch (...) { delete v; return nullptr; }
}

typedef int (*AudioCallback)(const float *, size_t, int32_t, void *);
extern "C" __declspec(dllexport) int voice_stream(void *handle, const char *text, AudioCallback emit, void *user) {
    auto v = static_cast<Voice *>(handle);
    try {
        p_mtmd_helper_gen_audio_reset(v->gen);
        p_llama_memory_clear(p_llama_get_memory(v->ctx), true);
        p_llama_sampler_reset(v->sampler);
        auto in = p_mtmd_helper_gen_audio_inp_default();
        in.stream = true; in.prompt = text; in.prompt_len = strlen(text);
        in.speaker_ref = v->speaker; in.lang = "en"; in.top_k = 40; in.top_p = 0.9f; in.seed = 42;
        in.out_type = MTMD_HELPER_GEN_AUDIO_OUTTYPE_PCM;
        if (p_mtmd_helper_gen_audio_set_input(v->gen, &in)) return 1;
        int prompt;
        do { prompt = p_mtmd_helper_gen_audio_step_prompt(v->gen, 512); } while (prompt > 0);
        if (prompt < 0) return 2;
        auto output = [&]() {
            int32_t rate = 0; const char *data = nullptr; size_t bytes = 0;
            if (p_mtmd_helper_gen_audio_get_output(v->gen, &rate, &data, &bytes, nullptr)) return false;
            return !bytes || emit(reinterpret_cast<const float *>(data), bytes/sizeof(float), rate, user) != 0;
        };
        auto sampled = p_llama_sampler_sample(v->sampler, v->ctx, -1);
        const float *hidden = p_llama_get_embeddings_ith(v->ctx, -1);
        bool stop = false;
        for (int n = 0; n < 768; ++n) {
            const float *next = nullptr;
            if (p_mtmd_helper_gen_audio_step_gen(v->gen, sampled, hidden, &next, &stop)) return 3;
            // The public flush API retains codec state across batches. Its
            // default 72-frame window delays short replies until completion.
            if ((n+1)%8 == 0 && p_mtmd_helper_gen_audio_flush(v->gen)) return 5;
            if (!output()) return 4;
            if (stop || !next) {
                if (p_mtmd_helper_gen_audio_flush(v->gen) || !output()) return 5;
                return 0;
            }
            hidden = next;
            sampled = p_llama_sampler_sample(v->sampler, v->ctx, -1);
        }
        return 6; // Never signal successful completion for truncated speech.
    } catch (...) { return 7; }
}
extern "C" __declspec(dllexport) void voice_free(void *handle) { delete static_cast<Voice *>(handle); }
