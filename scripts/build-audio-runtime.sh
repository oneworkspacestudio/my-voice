# Native Git Bash/MSYS make, Zig 0.15.2. Source digests: runtime.lock.json.
set -eu
ART=$(cd "${1:?Pass source/cache folder}" && pwd)
ROOT=$(cd "${2:?Pass native NTFS build folder}" && pwd)
ZIG=${3:?Pass Zig 0.15.2 executable}
export PATH="$ART/make/usr/bin:/usr/bin:/bin:$PATH"
export ZIG_GLOBAL_CACHE_DIR="$(cygpath -m "$ART/zig-cache")"
export ZIG_LOCAL_CACHE_DIR="$(cygpath -m "$ART/zig-local")"
mkdir -p "$ROOT" "$ROOT/prefix"
cd "$ROOT"
tar -xf "$ART/lame-source.tar.gz"
tar -xf "$ART/ffmpeg-source.tar.gz"
CC="$ZIG cc -target x86_64-windows-gnu"
AR="$ZIG ar"
RANLIB="$ZIG ranlib"
cd lame-3.100
if [ ! -f "$ROOT/prefix/lib/libmp3lame.a" ]; then
./configure --host=x86_64-w64-mingw32 --build=x86_64-pc-msys --prefix="$(cygpath -m "$ROOT/prefix")" --disable-shared --enable-static --disable-frontend CC="$CC" AR="$AR" RANLIB="$RANLIB" LD="$ZIG ld.lld" CFLAGS='-O2'
make -j6
make install
fi
cd "$ROOT"/FFmpeg-*
./configure --enable-cross-compile --target-os=mingw32 --arch=x86_64 --host-cc="$CC" --cc="$CC" --ar="$AR" --ranlib="$RANLIB" --disable-x86asm --disable-autodetect --disable-everything --disable-doc --enable-debug=0 --disable-stripping --disable-network --disable-avdevice --disable-ffplay --disable-ffprobe --enable-small --enable-ffmpeg --enable-protocol=file,pipe --enable-demuxer=wav,mp3,mov,matroska,ogg --enable-muxer=wav,mp3 --enable-decoder=pcm_s16le,pcm_s24le,pcm_s32le,pcm_f32le,pcm_f64le,mp3,mp3float,aac,opus,vorbis,flac --enable-encoder=pcm_s16le,libmp3lame --enable-parser=mpegaudio,aac,opus,vorbis,flac --enable-filter=aresample,asetrate,atempo,volume,anull,aformat,abuffer,abuffersink --enable-libmp3lame --extra-cflags="-I$(cygpath -m "$ROOT/prefix/include")" --extra-ldflags="-L$(cygpath -m "$ROOT/prefix/lib")" --prefix="$(cygpath -m "$ROOT/install")"
make -j6
cp ffmpeg.exe "$ART/ffmpeg-minimal.exe"
./ffmpeg.exe -version > "$ART/minimal-version.txt"
echo 'Native minimal audio runtime built successfully'
