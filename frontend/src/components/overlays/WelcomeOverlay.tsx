import React, { useState } from 'react';
import {
  Compass,
  Dices,
  DoorOpen,
  Loader2,
  LogIn,
  MapPinned,
  MessageCircle,
  Trophy,
  UserPlus,
} from 'lucide-react';
import { useViewStore } from '../../store/useViewStore';
import { useCityStore } from '../../store/useCityStore';
import { useUserStore } from '../../store/useUserStore';

const CITIES = ['北京', '上海', '西安', '成都', '杭州', '苏州', '南京', '武汉', '重庆', '广州', '厦门', '洛阳'];

const FEATURES = [
  { icon: MapPinned, label: '城市地标', text: '从真实城市与文化坐标开始探索' },
  { icon: MessageCircle, label: '人物对话', text: '和城市人物进行 AI 情境交流' },
  { icon: Trophy, label: '成就足迹', text: '收藏城市访问、打卡海报与徽章' },
];

export const WelcomeOverlay: React.FC = () => {
  const { enter, setView } = useViewStore();
  const cityCount = useCityStore(s => s.cities.length);
  const {
    username: currentUsername,
    nickname: currentNickname,
    register,
    login,
  } = useUserStore();

  const [authMode, setAuthMode] = useState<'register' | 'login'>('register');
  const [authUsername, setAuthUsername] = useState('');
  const [authPassword, setAuthPassword] = useState('');
  const [authNickname, setAuthNickname] = useState('');
  const [authLoading, setAuthLoading] = useState(false);
  const [authError, setAuthError] = useState('');
  const [showAuthForm, setShowAuthForm] = useState(!currentUsername);

  const handleEnter = (mode: 'FREE_EXPLORE' | 'GAME_DICE') => {
    enter();
    setView(mode);
  };

  const validateForm = () => {
    const normalizedUsername = authUsername.trim();
    if (!/^[A-Za-z0-9_]{3,32}$/.test(normalizedUsername)) {
      return '账号需为 3-32 位字母、数字或下划线';
    }
    if (authPassword.length < 6) {
      return '密码至少 6 位';
    }
    if (authMode === 'register' && authNickname.trim().length > 64) {
      return '昵称不能超过 64 字';
    }
    return '';
  };

  const handleAuth = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (authLoading) return;

    const nextError = validateForm();
    if (nextError) {
      setAuthError(nextError);
      return;
    }

    setAuthLoading(true);
    setAuthError('');
    try {
      if (authMode === 'register') {
        await register(authUsername.trim(), authPassword, authNickname.trim());
      } else {
        await login(authUsername.trim(), authPassword);
      }
      setAuthPassword('');
      handleEnter('FREE_EXPLORE');
    } catch (error: any) {
      setAuthError(error?.message || (authMode === 'register' ? '注册失败' : '登录失败'));
    } finally {
      setAuthLoading(false);
    }
  };

  return (
    <>
      <style>{`
        @keyframes welcomeFloat {
          0%, 100% { transform: translateY(0px); opacity: 0.12; }
          50% { transform: translateY(-18px); opacity: 0.24; }
        }
        @keyframes welcomeFadeUp {
          from { opacity: 0; transform: translateY(24px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .wf-1 { animation: welcomeFloat 7s ease-in-out infinite; }
        .wf-2 { animation: welcomeFloat 9s ease-in-out infinite 1.5s; }
        .wf-3 { animation: welcomeFloat 8s ease-in-out infinite 3s; }
        .wfu-0 { animation: welcomeFadeUp 0.65s ease-out 0.08s both; }
        .wfu-1 { animation: welcomeFadeUp 0.65s ease-out 0.18s both; }
        .wfu-2 { animation: welcomeFadeUp 0.65s ease-out 0.3s both; }
        .wfu-3 { animation: welcomeFadeUp 0.65s ease-out 0.42s both; }
      `}</style>

      <div className="fixed inset-0 z-50 overflow-hidden bg-background">
        <div
          className="absolute inset-0 opacity-80"
          style={{
            backgroundImage:
              'linear-gradient(rgba(43,58,54,0.055) 1px, transparent 1px), linear-gradient(90deg, rgba(43,58,54,0.055) 1px, transparent 1px)',
            backgroundSize: '56px 56px',
          }}
        />
        <div className="absolute inset-x-0 top-0 h-56 bg-gradient-to-b from-accent/15 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-72 bg-gradient-to-t from-primary/10 to-transparent" />

        <div className="absolute inset-0 pointer-events-none select-none overflow-hidden">
          {CITIES.map((city, i) => (
            <span
              key={city}
              className={`absolute text-primary font-serif-display font-bold text-sm wf-${(i % 3) + 1}`}
              style={{
                left: `${(i * 19 + 4) % 88}%`,
                top: `${(i * 17 + 6) % 82}%`,
                animationDelay: `${i * 0.4}s`,
                fontSize: `${12 + (i % 3) * 3}px`,
              }}
            >
              {city}
            </span>
          ))}
        </div>

        <main className="relative z-10 h-full overflow-y-auto">
          <div className="mx-auto grid min-h-full w-full max-w-6xl grid-cols-1 items-center gap-8 px-5 py-8 sm:px-8 lg:grid-cols-[minmax(0,1.08fr)_minmax(360px,440px)] lg:px-10">
            <section className="wfu-0 flex min-w-0 flex-col items-start text-left">
              <div className="mb-8 flex items-center gap-3">
                <img
                  src="/icon-transparent.png"
                  alt="任意门"
                  className="h-16 w-16 shrink-0 object-contain drop-shadow-[0_14px_28px_rgba(43,58,54,0.18)] sm:h-20 sm:w-20"
                />
                <div>
                  <div className="inline-flex items-center gap-2 rounded-full border border-accent/30 bg-accent/10 px-3 py-1 text-[10px] font-semibold tracking-[0.26em] text-accent">
                    AI 城市漫游
                  </div>
                  <div className="mt-2 text-sm font-semibold text-muted-foreground">推开门，遇见大美中国</div>
                </div>
              </div>

              <h1
                className="wfu-1 font-serif-display text-[64px] font-black leading-none text-foreground sm:text-[88px] lg:text-[104px]"
                style={{ textShadow: '0 16px 48px rgba(43,58,54,0.14), 0 0 80px rgba(194,159,96,0.18)' }}
              >
                任意门
              </h1>

              <p className="wfu-2 mt-5 max-w-xl text-base leading-8 text-foreground/68 sm:text-lg">
                随机降落中国城市，遇见历史地标，与古今人物跨时空对话，生成专属赛博打卡大片。
              </p>

              <div className="wfu-3 mt-10 hidden w-full max-w-2xl grid-cols-1 gap-3 sm:grid sm:grid-cols-3">
                {FEATURES.map(({ icon: Icon, label, text }) => (
                  <div key={label} className="rounded-lg border border-border/70 bg-card/70 p-4 shadow-sm backdrop-blur">
                    <div className="mb-4 flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <Icon className="h-5 w-5" />
                    </div>
                    <div className="text-sm font-bold text-foreground">{label}</div>
                    <div className="mt-1 text-xs leading-5 text-muted-foreground">{text}</div>
                  </div>
                ))}
              </div>

              <div className="wfu-3 mt-8 hidden flex-wrap gap-2 sm:flex">
                {[`${cityCount > 0 ? `${cityCount}座` : '多座'}城市`, 'AI人物对话', '赛博打卡', '成就收集'].map((tag) => (
                  <span
                    key={tag}
                    className="rounded-full border border-border/70 bg-background/80 px-3 py-1 text-xs font-medium text-muted-foreground shadow-sm"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </section>

            <aside className="wfu-2 w-full">
              <div className="rounded-lg border border-border/70 bg-background/95 p-5 shadow-[0_24px_60px_rgba(43,58,54,0.14)] backdrop-blur-xl sm:p-6">
                {currentUsername && !showAuthForm ? (
                  <div className="space-y-5">
                    <div>
                      <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                        <DoorOpen className="h-6 w-6" />
                      </div>
                      <h2 className="text-xl font-bold text-foreground">继续出发</h2>
                      <p className="mt-2 text-sm leading-6 text-muted-foreground">
                        {currentNickname || currentUsername} 的城市足迹已准备好。
                      </p>
                    </div>

                    <div className="grid grid-cols-1 gap-3">
                      <button
                        onClick={() => handleEnter('FREE_EXPLORE')}
                        className="flex h-12 items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm font-bold text-primary-foreground shadow-sm transition-colors hover:bg-primary/90"
                      >
                        <MapPinned className="h-5 w-5" />
                        自由探索
                      </button>
                      <button
                        onClick={() => handleEnter('GAME_DICE')}
                        className="flex h-12 items-center justify-center gap-2 rounded-lg border border-accent/40 bg-accent/10 px-4 text-sm font-bold text-primary transition-colors hover:bg-accent/20"
                      >
                        <Dices className="h-5 w-5" />
                        随机漫游
                      </button>
                    </div>

                    <button
                      onClick={() => {
                        setShowAuthForm(true);
                        setAuthError('');
                      }}
                      className="w-full rounded-lg border border-border px-4 py-2.5 text-sm font-semibold text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
                    >
                      使用其他账号
                    </button>
                  </div>
                ) : (
                  <div className="space-y-5">
                    <div>
                      <div className="mb-4 grid grid-cols-2 gap-1 rounded-lg bg-secondary p-1">
                        <button
                          type="button"
                          onClick={() => {
                            setAuthMode('register');
                            setAuthError('');
                          }}
                          className={`flex h-10 items-center justify-center gap-2 rounded-md text-sm font-bold transition-colors ${
                            authMode === 'register'
                              ? 'bg-background text-foreground shadow-sm'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          <UserPlus className="h-4 w-4" />
                          注册
                        </button>
                        <button
                          type="button"
                          onClick={() => {
                            setAuthMode('login');
                            setAuthError('');
                          }}
                          className={`flex h-10 items-center justify-center gap-2 rounded-md text-sm font-bold transition-colors ${
                            authMode === 'login'
                              ? 'bg-background text-foreground shadow-sm'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          <LogIn className="h-4 w-4" />
                          登录
                        </button>
                      </div>
                      <h2 className="text-xl font-bold text-foreground">
                        {authMode === 'register' ? '创建城市漫游账号' : '登录城市漫游账号'}
                      </h2>
                      <p className="mt-2 text-sm leading-6 text-muted-foreground">
                        {authMode === 'register' ? '注册后保留匿名足迹、成就与打卡海报。' : '登录后同步你的城市足迹。'}
                      </p>
                    </div>

                    <form onSubmit={handleAuth} className="space-y-3">
                      <label className="block">
                        <span className="mb-1.5 block text-xs font-semibold text-muted-foreground">账号</span>
                        <input
                          value={authUsername}
                          onChange={(event) => setAuthUsername(event.target.value)}
                          placeholder="demo_user"
                          autoComplete="username"
                          className="h-11 w-full rounded-lg border border-input bg-card px-3 text-sm outline-none transition-colors placeholder:text-muted-foreground/50 focus:border-primary/50 focus:bg-background"
                        />
                      </label>
                      <label className="block">
                        <span className="mb-1.5 block text-xs font-semibold text-muted-foreground">密码</span>
                        <input
                          type="password"
                          value={authPassword}
                          onChange={(event) => setAuthPassword(event.target.value)}
                          placeholder="至少 6 位"
                          autoComplete={authMode === 'register' ? 'new-password' : 'current-password'}
                          className="h-11 w-full rounded-lg border border-input bg-card px-3 text-sm outline-none transition-colors placeholder:text-muted-foreground/50 focus:border-primary/50 focus:bg-background"
                        />
                      </label>
                      {authMode === 'register' && (
                        <label className="block">
                          <span className="mb-1.5 block text-xs font-semibold text-muted-foreground">昵称</span>
                          <input
                            value={authNickname}
                            onChange={(event) => setAuthNickname(event.target.value)}
                            placeholder="北京游客"
                            autoComplete="nickname"
                            className="h-11 w-full rounded-lg border border-input bg-card px-3 text-sm outline-none transition-colors placeholder:text-muted-foreground/50 focus:border-primary/50 focus:bg-background"
                          />
                        </label>
                      )}

                      {authError && (
                        <div className="rounded-lg border border-destructive/20 bg-destructive/10 px-3 py-2 text-xs font-medium text-destructive">
                          {authError}
                        </div>
                      )}

                      <button
                        type="submit"
                        disabled={authLoading}
                        className="flex h-12 w-full items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm font-bold text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
                      >
                        {authLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : authMode === 'register' ? <UserPlus className="h-4 w-4" /> : <LogIn className="h-4 w-4" />}
                        {authLoading ? '处理中...' : authMode === 'register' ? '注册并进入' : '登录并进入'}
                      </button>
                    </form>

                    <div className="grid grid-cols-2 gap-2 border-t border-border/70 pt-4">
                      <button
                        type="button"
                        onClick={() => handleEnter('FREE_EXPLORE')}
                        className="flex h-11 items-center justify-center gap-2 rounded-lg border border-border bg-background px-3 text-xs font-bold text-foreground transition-colors hover:bg-secondary"
                      >
                        <Compass className="h-4 w-4" />
                        游客探索
                      </button>
                      <button
                        type="button"
                        onClick={() => handleEnter('GAME_DICE')}
                        className="flex h-11 items-center justify-center gap-2 rounded-lg border border-accent/40 bg-accent/10 px-3 text-xs font-bold text-primary transition-colors hover:bg-accent/20"
                      >
                        <Dices className="h-4 w-4" />
                        随机漫游
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </aside>
          </div>
        </main>
      </div>
    </>
  );
};
