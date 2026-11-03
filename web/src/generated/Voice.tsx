import { Alert, Badge, Button, Card, Checkbox, Divider, Grid, Heading, Icon, Image, Input, Link, Page, Row, Section, Select, Spacer, Stack, Text } from "@/components/ui"

export function Voice({ data, onAction }: { data: any; onAction: (name: string) => void }) {
  return (
    <Page>
      {data.loggedOut ? (
        <Card variant="login">
          <Image src="./logo.svg" alt="My Voice" width={56} height={56} />
          <Heading level={1}>Your words. Your voice.</Heading>
          <Text variant="muted">Sign in to your private audio workspace.</Text>
          {data.oidc ? (
            <Link href="./auth/login" variant="primary">Continue with MyAccounts</Link>
          ) : (
            <>
              <Input name="email" label="Email" type="email" />
              <Input name="password" label="Password" type="password" />
              <Button onClick={() => onAction("login")} variant="primary">{data.signInLabel}</Button>
            </>
          )}
          {data.error && (
            <Alert>{data.error}</Alert>
          )}
        </Card>
      ) : (
        <>
          <Stack variant="sidebar">
            <Row variant="brand">
              <Image src="./logo.svg" alt="" width={34} height={34} />
              <Heading level={1}>My Voice</Heading>
            </Row>
            <Text variant="eyebrow">YOUR AUDIO WORKSPACE</Text>
            {data.nav.map((item, i) => (
              <Button key={i} onClick={() => onAction("navigate")} item={item.id} icon={item.icon} label={item.label} variant={item.variant} />
            ))}
            <Divider />
            <Text variant="eyebrow">CREATE</Text>
            <Button onClick={() => onAction("speech")} icon="AudioLines" label="Text to Speech" variant="nav" />
            <Button onClick={() => onAction("clone")} icon="Fingerprint" label="Voice Cloning" variant="nav" />
            <Button onClick={() => onAction("dialogue")} icon="MessagesSquare" label="Multi-voice dialogue" variant="nav" />
            <Button onClick={() => onAction("batch")} icon="Layers" label="Batch generation" variant="nav" />
            <Spacer />
            <Card variant="privacy">
              <Icon name="ShieldCheck" size={21} />
              <Text variant="privacy-title">Made to stay yours</Text>
              <Text variant="caption">Local speech. Private voices. No subscription.</Text>
            </Card>
            <Link href={data.workspace} variant="apps">
              <Icon name="Grid2X2" size={17} />
              <Text>One Workspace</Text>
            </Link>
            <Text variant="version">{data.version}</Text>
          </Stack>
          <Stack variant="main">
            <Row variant="announcement">
              <Icon name="Sparkles" size={14} />
              <Text>A little text. A world of possibilities.</Text>
              <Badge>CPU powered</Badge>
            </Row>
            <Row variant="header">
              <Button onClick={() => onAction("menu")} icon="Menu" label="Open navigation" variant="mobile-menu" />
              <Icon name={data.pageIcon} size={19} />
              <Text>{data.pageTitle}</Text>
              <Spacer />
              <Input name="search" placeholder="Search your creations…" variant="search" />
              <Link href={data.workspace} variant="apps">
                <Icon name="Grid2X2" size={16} />
                <Text>Apps</Text>
              </Link>
              <Button onClick={() => onAction("help")} icon="CircleHelp" label="Help" variant="icon" />
              <Button onClick={() => onAction("account")} variant="avatar">{data.initials}</Button>
            </Row>
            {data.error && (
              <Row variant="notice">
                <Alert>{data.error}</Alert>
                <Button onClick={() => onAction("dismiss")} icon="X" label="Dismiss message" variant="icon" />
              </Row>
            )}
            {data.hasPreview && (
              <Card variant="voice-preview">
                <Row>
                  <Icon name="AudioLines" size={21} />
                  <Stack>
                    <Heading level={2}>{data.previewName}</Heading>
                    <Text variant="caption">{data.previewStatus}</Text>
                  </Stack>
                  <Spacer />
                  {data.previewRunning && (
                    <Button onClick={() => onAction("cancel-preview")} icon="Square" label="Stop preview" variant="secondary" />
                  )}
                  {data.previewFailed && (
                    <Button onClick={() => onAction("retry-preview")} icon="RefreshCw" label="Retry preview" variant="secondary" />
                  )}
                  <Button onClick={() => onAction("close-preview")} icon="X" label="Close voice preview" variant="icon" />
                </Row>
                <Section id="preview-audio" />
              </Card>
            )}
            {data.home && (
              <Stack variant="home">
                <Stack variant="hero">
                  <Badge>MY VOICE · YOUR CREATIVE SPACE</Badge>
                  <Heading level={1}>What would you like to create?</Heading>
                  <Text variant="muted">Bring your words to life. Make them sound like you.</Text>
                  <Row variant="prompt">
                    <Icon name="AudioLines" size={26} />
                    <Input name="idea" placeholder="Start with a story, an idea, a few words…" variant="idea" />
                    <Button onClick={() => onAction("start")} icon="ArrowUp" label="Start creating" variant="round" />
                  </Row>
                </Stack>
                <Row variant="quick-tools">
                  {data.quickTools.map((item, i) => (
                    <Button key={i} onClick={() => onAction("navigate")} item={item.id} icon={item.icon} label={item.label} variant="quick-tool" />
                  ))}
                </Row>
                <Card variant="feature">
                  <Stack>
                    <Badge>YOUR NEXT GREAT IDEA, SPOKEN</Badge>
                    <Heading level={2}>A voice for every story.</Heading>
                    <Text>Natural speech and instant voice cloning, right in your Workspace.</Text>
                    <Row>
                      <Button onClick={() => onAction("speech")} icon="ArrowUpRight" label="Create your first speech" variant="white" />
                      <Button onClick={() => onAction("clone")} label="Make it your voice" variant="feature-link" />
                    </Row>
                  </Stack>
                  <Section id="feature-wave" />
                </Card>
                <Row variant="section-heading">
                  <Heading level={2}>Make something great</Heading>
                  <Spacer />
                  <Button onClick={() => onAction("history")} variant="text">View all creations</Button>
                </Row>
                <Grid variant="quickstarts">
                  {data.starts.map((item, i) => (
                    <Button key={i} onClick={() => onAction("navigate")} item={item.id} variant="start-card" icon={item.icon} label={item.label} detail={item.detail} />
                  ))}
                </Grid>
                <Row variant="section-heading">
                  <Heading level={2}>Recent creations</Heading>
                  <Spacer />
                  <Text variant="caption">Saved privately to your account</Text>
                </Row>
                {data.hasHistory ? (
                  <Stack variant="recent">
                    {data.recent.map((item, i) => (
                      <Button key={i} onClick={() => onAction("open-job")} item={item.id} variant="recent" icon="AudioLines" label={item.script.name} detail={item.detail} />
                    ))}
                  </Stack>
                ) : (
                  <Card variant="empty-recent">
                    <Icon name="AudioLines" size={23} />
                    <Text>Your next creation starts with a few words.</Text>
                    <Button onClick={() => onAction("speech")} variant="text">Write something</Button>
                  </Card>
                )}
              </Stack>
            )}
            {data.speech && (
              <Stack variant="studio">
                <Row variant="titlebar">
                  <Stack>
                    <Heading level={1}>Text to Speech</Heading>
                    <Text variant="muted">A blank page. Endless ways to be heard.</Text>
                  </Stack>
                  <Spacer />
                  <Button onClick={() => onAction("import")} icon="Upload" label="Import text" variant="secondary" />
                  <Button onClick={() => onAction("import-drive")} icon="CloudDownload" label="My Drive" variant="secondary" />
                  <Button onClick={() => onAction("new-project")} icon="Plus" label="New project" variant="secondary" />
                </Row>
                <Row variant="studio-layout">
                  <Stack variant="editor-column">
                    <Card variant="editor">
                      <Row variant="editor-heading">
                        <Input name="projectName" placeholder="Untitled speech" variant="title" />
                        <Badge>{data.saved}</Badge>
                      </Row>
                      <Row variant="editor-tools">
                        <Button onClick={() => onAction("undo")} icon="Undo2" label="Undo" variant="icon" />
                        <Button onClick={() => onAction("redo")} icon="Redo2" label="Redo" variant="icon" />
                        <Button onClick={() => onAction("dictionary")} icon="BookOpen" label="Pronunciation" variant="text" />
                        <Button onClick={() => onAction("toggle-dialogue")} icon="MessagesSquare" label={data.segmentDetailsLabel} variant={data.dialogueVariant} expanded={data.showDialogue} controls="speech-segments" />
                        <Spacer />
                        <Text variant="caption">TXT · MD</Text>
                      </Row>
                      <Input name="script" label="Speech script" multiline={true} placeholder="Write or paste your script here. A story, a welcome message, the beginning of something new…" variant="script" />
                      <Row variant="editor-footer">
                        <Text variant="caption">{data.counter}</Text>
                        <Spacer />
                        <Text variant="caption">Ctrl + Enter to generate</Text>
                      </Row>
                    </Card>
                    {data.showDialogue && (
                      <Card variant="segments" id="speech-segments">
                        <Row>
                          <Heading level={2}>Speakers & segments</Heading>
                          <Spacer />
                          <Button onClick={() => onAction("split")} icon="Scissors" label="Split script" variant="text" />
                        </Row>
                        <Text variant="caption">Assign each segment a voice. After generating, edit and regenerate just one segment.</Text>
                        {data.segments.map((item, i) => (
                          <Row key={i} variant="segment">
                            <Text variant="segment-number">{item.number}</Text>
                            <Input name={item.textField} label={item.label} multiline={true} variant="segment-text" />
                            <Select name={item.voiceField} label="Segment voice" options={data.voiceOptions} />
                            <Button onClick={() => onAction("regenerate")} item={item.id} icon="RefreshCw" label="Regenerate segment" variant="icon" disabled={item.disabled} />
                          </Row>
                        ))}
                      </Card>
                    )}
                    {data.showDictionary && (
                      <Card>
                        <Heading level={2}>Pronunciation dictionary</Heading>
                        <Text variant="caption">Exact text replacements, applied before speech. Project rules override account rules; no SSML or phoneme support is implied.</Text>
                        <Input name="rules" label="Project replacements" multiline={true} placeholder="One Workspace = One Work space" variant="dictionary" />
                        <Text variant="caption">{data.defaultRuleText}</Text>
                        <Row>
                          <Button onClick={() => onAction("preview-text")} icon="Eye" label="Preview spoken text" variant="secondary" />
                          <Button onClick={() => onAction("save-dictionary")} label="Save as account defaults" variant="text" />
                        </Row>
                        <Text variant="pre">{data.transformed}</Text>
                      </Card>
                    )}
                    {data.hasJob && (
                      <Card variant="player">
                        <Row>
                          <Icon name="AudioLines" size={23} />
                          <Stack>
                            <Heading level={2}>{data.jobName}</Heading>
                            <Text variant="caption">{data.jobStatus}</Text>
                          </Stack>
                          <Spacer />
                          {data.jobRunning && (
                            <Button onClick={() => onAction("cancel")} icon="Square" label="Cancel" variant="secondary" />
                          )}
                          {data.jobFailed && (
                            <Button onClick={() => onAction("retry")} icon="RefreshCw" label="Retry" variant="secondary" />
                          )}
                        </Row>
                        <Section id="waveform" />
                        <Section id="audio" />
                        {data.jobComplete && (
                          <Row variant="exports">
                            <Link href={data.wav} variant="secondary">Download WAV</Link>
                            <Link href={data.mp3} variant="secondary">Download MP3</Link>
                            <Button onClick={() => onAction("drive")} icon="CloudUpload" label="Save to My Drive" variant="text" />
                          </Row>
                        )}
                      </Card>
                    )}
                    {data.showBatch && (
                      <Card>
                        <Heading level={2}>Batch generation</Heading>
                        <Text variant="caption">One recording per line. Up to 20 queued jobs, processed one at a time.</Text>
                        <Input name="batchText" label="Batch scripts" multiline={true} variant="dictionary" />
                        <Row>
                          <Button onClick={() => onAction("run-batch")} icon="Layers" label="Generate batch" variant="primary" />
                          <Button onClick={() => onAction("export-batch")} icon="Download" label="Export completed ZIP" variant="secondary" />
                        </Row>
                      </Card>
                    )}
                  </Stack>
                  <Card variant="voice-settings">
                    <Row variant="tabs">
                      <Text variant="tab-active">Settings</Text>
                      <Spacer />
                      <Text variant="caption">Local generation</Text>
                    </Row>
                    <Text variant="field-label">Voice</Text>
                    <Select name="voice" label="Voice" options={data.voiceOptions} />
                    <Row variant="selected-voice">
                      <Icon name="AudioLines" size={19} />
                      <Stack>
                        <Text>{data.voiceName}</Text>
                        <Text variant="caption">{data.voiceDescription}</Text>
                      </Stack>
                      <Spacer />
                      <Button onClick={() => onAction("preview-voice")} icon="Play" label="Preview voice" variant="icon" />
                    </Row>
                    <Button onClick={() => onAction("clone")} icon="Plus" label="Create a voice" variant="text" />
                    <Divider />
                    <Select name="model" label="Model" options={data.modelOptions} />
                    <Text variant="caption">{data.modelDescription}</Text>
                    <Select name="language" label="Language" options={data.languages} />
                    <Divider />
                    <Row>
                      <Text variant="field-label">Audio controls</Text>
                      <Badge>Post-processing</Badge>
                    </Row>
                    <Input name="speed" label={data.speedLabel} type="range" min={0.5} max={2} step={0.05} />
                    <Input name="pitch" label={data.pitchLabel} type="range" min={-12} max={12} step={1} />
                    <Input name="volume" label={data.volumeLabel} type="range" min={0.1} max={2} step={0.05} />
                    <Input name="pause" label={data.pauseLabel} type="range" min={0} max={2} step={0.1} />
                    <Text variant="caption">Speed, pitch and volume are applied to generated audio. Voice identity stays fixed to your selected model.</Text>
                    <Spacer />
                    <Button onClick={() => onAction("generate")} icon="Sparkles" label="Generate speech" variant="primary" disabled={data.generating} />
                    <Text variant="caption">Private by default. Generated on this PC.</Text>
                  </Card>
                </Row>
              </Stack>
            )}
            {data.clone && (
              <Stack variant="clone-page">
                <Row variant="titlebar">
                  <Stack>
                    <Heading level={1}>Make it your voice.</Heading>
                    <Text variant="muted">A few seconds of you. A whole world of possibilities.</Text>
                  </Stack>
                  <Spacer />
                  <Badge>INSTANT VOICE CLONING</Badge>
                </Row>
                <Row variant="steps">
                  <Badge>1 · Add a reference</Badge>
                  <Icon name="ChevronRight" size={16} />
                  <Text>2 · Test your voice</Text>
                  <Icon name="ChevronRight" size={16} />
                  <Text>3 · Ready to create</Text>
                </Row>
                <Card variant="clone-model-setup">
                  <Row>
                    <Icon name="Cpu" size={28} />
                    <Stack>
                      <Heading level={2}>{data.cloneSetupTitle}</Heading>
                      <Text variant="muted">{data.cloneSetupDetail}</Text>
                    </Stack>
                    <Spacer />
                    <Badge>{data.cloneSetupStatus}</Badge>
                  </Row>
                  <Select name="cloneModel" label="Cloning model" options={data.cloneOptions} disabled={data.cloneModelDisabled} />
                  <Section id="clone-model-progress" />
                  <Row variant="clone-model-actions">
                    {data.cloneCanInstall && (
                      <Button onClick={() => onAction("install-cloning-model")} icon="Download" label={data.cloneInstallLabel} variant="primary" />
                    )}
                    {data.cloneNeedsOwner && (
                      <Text variant="caption">Ask the Workspace owner to install this model. You can record once it is ready.</Text>
                    )}
                    <Button onClick={() => onAction("refresh-models")} icon="RefreshCw" label="Refresh models" variant="text" />
                  </Row>
                </Card>
                <Card variant="clone-card">
                  <Row>
                    <Icon name="Fingerprint" size={32} />
                    <Stack>
                      <Heading level={2}>Every great voice starts with a sample.</Heading>
                      <Text variant="muted">Speak naturally. We trim silence, balance volume and choose a 3–30 second reference automatically.</Text>
                    </Stack>
                  </Row>
                  <Text variant="caption">For longer uploads, put the speech you want to use in the first two minutes.</Text>
                  <Section id="upload" />
                  <Row variant="record-actions">
                    <Button onClick={() => onAction("record")} icon="Mic" label={data.recordLabel} variant="secondary" disabled={data.cloneRecordDisabled} />
                    <Text variant="caption">WAV, MP3, M4A · up to 25 MB</Text>
                  </Row>
                  <Section id="reference-audio" />
                  <Grid variant="two">
                    <Input name="cloneName" label="Voice name" placeholder="Give your voice a name" />
                    <Select name="cloneLanguage" label="Reference language" options={data.cloneLanguages} />
                  </Grid>
                  <Input name="cloneTags" label="Tags" placeholder="Warm, storytelling, my voice…" />
                  <Checkbox name="consent" label="I own this voice, or have explicit permission from the speaker to clone and use it." />
                  <Text variant="caption">Your reference and model-specific voice state stay private. Share use only with people you explicitly choose.</Text>
                  <Row variant="clone-submit">
                    <Button onClick={() => onAction("save-clone")} icon="Fingerprint" label="Create my voice" variant="primary" disabled={data.cloneCreateDisabled} />
                    <Text variant="caption">{data.cloneStatus}</Text>
                  </Row>
                  {data.cloneCreated && (
                    <Card variant="clone-success">
                      <Icon name="CircleCheck" size={28} />
                      <Heading level={2}>Your voice is ready.</Heading>
                      <Text>Your original recording is kept. Listen to the prepared reference, then compare it with your generated voice.</Text>
                      <Text variant="caption">Prepared reference</Text>
                      <Section id="prepared-reference-audio" />
                      <Button onClick={() => onAction("test-clone")} icon="Play" label="Generate a test sample" variant="secondary" />
                      <Section id="clone-test-audio" />
                      <Button onClick={() => onAction("use-clone")} icon="ArrowRight" label="Use this voice" variant="primary" />
                    </Card>
                  )}
                </Card>
                <Card variant="tips">
                  <Heading level={2}>A little preparation makes a big difference.</Heading>
                  <Grid variant="quickstarts">
                    <Stack>
                      <Icon name="AudioLines" size={22} />
                      <Text>Keep it natural</Text>
                      <Text variant="caption">Speak at your normal pace, with a little expression.</Text>
                    </Stack>
                    <Stack>
                      <Icon name="VolumeX" size={22} />
                      <Text>Find a quiet spot</Text>
                      <Text variant="caption">Avoid music, echo and background conversations.</Text>
                    </Stack>
                    <Stack>
                      <Icon name="ShieldCheck" size={22} />
                      <Text>Keep it permissioned</Text>
                      <Text variant="caption">Only clone voices you own or are authorised to use.</Text>
                    </Stack>
                  </Grid>
                </Card>
              </Stack>
            )}
            {data.voicesPage && (
              <Stack variant="library-page">
                <Row variant="titlebar">
                  <Stack>
                    <Heading level={1}>A library that sounds like you.</Heading>
                    <Text variant="muted">Your voices, ready for your next idea.</Text>
                  </Stack>
                  <Spacer />
                  <Button onClick={() => onAction("clone")} icon="Plus" label="Create a voice" variant="primary" />
                </Row>
                <Row variant="filters">
                  <Input name="voiceSearch" placeholder="Search voices or tags…" />
                  <Select name="voiceFilter" label="Voice collection" options={data.voiceFilters} />
                  <Select name="voiceLanguage" label="Reference language filter" options={data.voiceLanguages} />
                </Row>
                <Grid variant="voices">
                  {data.filteredVoices.map((item, i) => (
                    <Card key={i} variant="voice-card">
                      <Row>
                        <Icon name="AudioLines" size={28} />
                        <Spacer />
                        <Badge>{item.type}</Badge>
                      </Row>
                      <Heading level={2}>{item.name}</Heading>
                      <Text variant="caption">{item.tags}</Text>
                      <Badge>{item.engineLabel}</Badge>
                      <Row>
                        <Button onClick={() => onAction("choose-voice")} item={item.id} icon="ArrowRight" label="Use voice" variant="secondary" />
                        <Button onClick={() => onAction("listen-voice")} item={item.id} icon="Play" label="Preview" variant="icon" />
                        <Button onClick={() => onAction("favorite-voice")} item={item.id} icon="Star" label={item.favoriteLabel} variant={item.favoriteVariant} />
                        {item.owned && (
                          <Button onClick={() => onAction("edit-voice")} item={item.id} icon="Ellipsis" label="Manage voice" variant="icon" />
                        )}
                      </Row>
                    </Card>
                  ))}
                </Grid>
              </Stack>
            )}
            {data.historyPage && (
              <Stack variant="library-page">
                <Row variant="titlebar">
                  <Stack>
                    <Heading level={1}>Every word, kept close.</Heading>
                    <Text variant="muted">Pick up where you left off. Listen, export or create again.</Text>
                  </Stack>
                  <Spacer />
                  <Button onClick={() => onAction("export-all")} icon="Download" label="Export completed ZIP" variant="secondary" />
                </Row>
                <Input name="historySearch" placeholder="Search recordings and scripts…" />
                {data.hasHistory ? (
                  <Stack>
                    {data.filteredHistory.map((item, i) => (
                      <Card key={i} variant="history-row">
                        <Icon name="AudioLines" size={24} />
                        <Stack>
                          <Heading level={2}>{item.script.name}</Heading>
                          <Text variant="caption">{item.detail}</Text>
                        </Stack>
                        <Spacer />
                        <Badge>{item.status}</Badge>
                        <Button onClick={() => onAction("open-job")} item={item.id} icon="Play" label="Open recording" variant="secondary" />
                        <Button onClick={() => onAction("delete-recording")} item={item.id} icon="Trash2" label="Remove recording" variant="icon" />
                      </Card>
                    ))}
                  </Stack>
                ) : (
                  <Card variant="empty">
                    <Icon name="AudioLines" size={42} />
                    <Heading level={2}>Your story starts here.</Heading>
                    <Text variant="muted">Create your first recording and it will be waiting here.</Text>
                    <Button onClick={() => onAction("speech")} label="Create speech" variant="primary" />
                  </Card>
                )}
              </Stack>
            )}
            {data.settingsPage && (
              <Stack variant="library-page">
                <Row variant="titlebar">
                  <Stack>
                    <Heading level={1}>Your voice. Your setup.</Heading>
                    <Text variant="muted">Choose the right balance for your PC.</Text>
                  </Stack>
                </Row>
                {data.models.map((item, i) => (
                  <Card key={i} variant="model-card">
                    <Row>
                      <Icon name="Cpu" size={28} />
                      <Stack>
                        <Heading level={2}>{item.name}</Heading>
                        <Text variant="caption">{item.runtime}</Text>
                      </Stack>
                      <Spacer />
                      <Badge>{item.installLabel}</Badge>
                    </Row>
                    <Text variant="caption">{item.license}</Text>
                    <Text variant="caption">{item.description}</Text>
                    <Row>
                      <Link href={item.url} variant="text">Model & licence details</Link>
                      <Spacer />
                      <Button onClick={() => onAction("install-model")} item={item.id} icon="Download" label={item.buttonLabel} variant="secondary" />
                      {item.installed && (
                        <Button onClick={() => onAction("remove-model")} item={item.id} icon="Trash2" label="Remove model" variant="icon" />
                      )}
                    </Row>
                    {item.status && (
                      <Text variant="caption">{item.status}</Text>
                    )}
                  </Card>
                ))}
                <Card>
                  <Heading level={2}>Performance & privacy</Heading>
                  <Input name="threads" label="CPU threads (1–8)" type="number" min={1} max={8} />
                  <Button onClick={() => onAction("save-settings")} label="Save settings" variant="primary" />
                  <Text variant="caption">{data.storageInfo}</Text>
                  <Text variant="caption">{data.retentionInfo}</Text>
                  <Text variant="caption">Local CPU generation. No audio sent to a speech API. Model downloads are explicit; installed models work offline.</Text>
                  <Text variant="caption">Voice references and clone states are private app data. Export your audio to My Drive or download a copy.</Text>
                </Card>
              </Stack>
            )}
            {data.modal && (
              <Card variant="modal">
                <Row>
                  <Heading level={2}>{data.modalTitle}</Heading>
                  <Spacer />
                  <Button onClick={() => onAction("close-modal")} icon="X" label="Close dialog" variant="icon" />
                </Row>
                {data.installModal && (
                  <>
                    <Text>{data.installDescription}</Text>
                    <Link href={data.installURL} variant="text">Read the model licence and access terms</Link>
                    <Checkbox name="license" label="I have reviewed and accept the model licence and usage conditions." />
                    {data.pocketInstall && (
                      <>
                        <Checkbox name="enableCloning" label="Enable Pocket voice cloning (requires authorised Hugging Face access)." />
                        <Input name="hfToken" type="password" label="Hugging Face read token (sent only to Hugging Face; never saved)" placeholder="hf_…" />
                      </>
                    )}
                    <Button onClick={() => onAction("confirm-install")} icon="Download" label="Download model" variant="primary" disabled={data.installDisabled} />
                  </>
                )}
                {data.manageModal && (
                  <>
                    <Input name="editName" label="Voice name" />
                    <Input name="editTags" label="Tags" />
                    <Checkbox name="favorite" label="Favourite voice" />
                    <Button onClick={() => onAction("save-voice")} label="Save voice" variant="primary" />
                    <Divider />
                    <Input name="shareEmail" label="Share use with OneAccount email" type="email" />
                    <Button onClick={() => onAction("share-voice")} icon="UserPlus" label="Grant private use" variant="secondary" />
                    <Text variant="caption">{data.shareResult}</Text>
                    {data.shares.map((item, i) => (
                      <Row key={i}>
                        <Text>{item.email}</Text>
                        <Spacer />
                        <Button onClick={() => onAction("revoke-specific")} item={item.user} icon="UserMinus" label="Revoke use" variant="secondary" />
                      </Row>
                    ))}
                    <Divider />
                    <Checkbox name="deleteConsent" label="Permanently remove this voice and its reference recording." />
                    <Button onClick={() => onAction("delete-voice")} icon="Trash2" label="Delete voice" variant="danger" />
                  </>
                )}
                {data.helpModal && (
                  <>
                    <Text>Start with Text to Speech. Install a CPU model in Settings, choose a voice and write your script. Use Voice Cloning only with a speaker’s permission.</Text>
                    <Text>Dialogue splits your script into segments. After generation, change a segment and press its regenerate button. TXT and Markdown import preserve sentence order.</Text>
                    <Text>My Voice runs through OneAccount and Workspace. For standalone use, create a local account using the one-time setup link printed by the launcher.</Text>
                  </>
                )}
                {data.accountModal && (
                  <>
                    <Text>{data.accountName}</Text>
                    <Button onClick={() => onAction("logout")} icon="LogOut" label="Sign out" variant="secondary" />
                  </>
                )}
                {data.removeModelModal && (
                  <>
                    <Text>Remove the downloaded model files? Your voices, scripts and generated recordings will be kept. You can reinstall the model later.</Text>
                    <Button onClick={() => onAction("confirm-remove-model")} icon="Trash2" label="Remove downloaded model" variant="danger" />
                  </>
                )}
                {data.deleteRecordingModal && (
                  <>
                    <Text>Permanently remove this recording and its generated segments? Your voice library and saved script projects will be kept.</Text>
                    <Checkbox name="deleteRecordingConsent" label="I want to permanently delete this recording." />
                    <Button onClick={() => onAction("confirm-delete-recording")} icon="Trash2" label="Permanently remove recording" variant="danger" />
                  </>
                )}
                {data.driveImportModal && (
                  <>
                    <Text>Choose a TXT or Markdown file from your signed-in My Drive account.</Text>
                    <Select name="driveSource" label="Drive source" options={data.driveSources} />
                    <Input name="drivePath" label="Folder path" />
                    <Button onClick={() => onAction("list-drive")} icon="FolderOpen" label="Open folder" variant="secondary" />
                    {data.driveFiles.map((item, i) => (
                      <Button key={i} onClick={() => onAction("select-drive-file")} item={item.path} icon={item.icon} label={item.name} variant="drive-file" />
                    ))}
                  </>
                )}
              </Card>
            )}
          </Stack>
        </>
      )}
    </Page>
  )
}
