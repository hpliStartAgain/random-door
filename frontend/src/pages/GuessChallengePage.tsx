import React, { useEffect, useState } from 'react';
import { ArrowLeft, CheckCircle2, XCircle } from 'lucide-react';
import { api } from '../api';
import type { GuessAnswerResponse, GuessChallengeResponse } from '../api/types';

interface Props {
  code: string;
}

export const GuessChallengePage: React.FC<Props> = ({ code }) => {
  const [challenge, setChallenge] = useState<GuessChallengeResponse | null>(null);
  const [answer, setAnswer] = useState('');
  const [result, setResult] = useState<GuessAnswerResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    setLoading(true);
    api.getGuessChallenge(code)
      .then(setChallenge)
      .catch((e: any) => setError(e?.message || '挑战不存在或已过期'))
      .finally(() => setLoading(false));
  }, [code]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const value = answer.trim();
    if (!value || submitting) return;
    setSubmitting(true);
    setError('');
    try {
      const res = await api.answerGuessChallenge(code, value);
      setResult(res);
    } catch (e: any) {
      setError(e?.message || '提交失败');
    } finally {
      setSubmitting(false);
    }
  };

  const backHome = () => {
    window.history.replaceState({}, '', window.location.pathname);
    window.location.reload();
  };

  return (
    <div className="min-h-screen bg-background text-foreground">
      <div className="max-w-3xl mx-auto px-4 py-8">
        <button
          onClick={backHome}
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors mb-6"
        >
          <ArrowLeft className="h-4 w-4" />
          返回任意门
        </button>

        <div className="rounded-2xl border border-border bg-card shadow-sm overflow-hidden">
          <div className="px-5 py-4 border-b border-border/50">
            <h1 className="font-serif-display text-2xl font-bold">猜猜我在哪</h1>
            <p className="text-sm text-muted-foreground mt-1">
              {challenge?.target_name || '一段城市风光'} · 好友发来的挑战
            </p>
          </div>

          {loading && <div className="p-8 text-sm text-muted-foreground">正在加载挑战...</div>}

          {!loading && error && !challenge && (
            <div className="p-8 text-sm text-red-500">{error}</div>
          )}

          {challenge && (
            <div className="p-5 space-y-5">
              {challenge.image_url && (
                <div className="rounded-xl overflow-hidden border border-border bg-background aspect-video">
                  <img src={challenge.image_url} alt="猜位置截图" className="w-full h-full object-cover" />
                </div>
              )}

              {challenge.caption && (
                <div className="rounded-xl border border-border bg-background/70 p-4 text-sm leading-relaxed text-foreground/80">
                  {challenge.caption}
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-3">
                <label className="text-xs font-semibold text-muted-foreground">你的答案</label>
                <div className="flex gap-2">
                  <input
                    value={answer}
                    onChange={(e) => setAnswer(e.target.value)}
                    disabled={!!result}
                    placeholder="输入城市名或地标名"
                    className="flex-1 px-4 py-3 rounded-xl border border-border bg-background text-sm outline-none focus:border-primary/50 disabled:opacity-60"
                  />
                  <button
                    type="submit"
                    disabled={!answer.trim() || submitting || !!result}
                    className="px-5 py-3 rounded-xl bg-primary text-primary-foreground text-sm font-bold disabled:opacity-50"
                  >
                    {submitting ? '提交中' : '提交'}
                  </button>
                </div>
              </form>

              {error && <div className="text-sm text-red-500">{error}</div>}

              {result && (
                <div className={`rounded-xl border p-4 flex gap-3 ${result.is_correct ? 'border-green-200 bg-green-50 text-green-900' : 'border-amber-200 bg-amber-50 text-amber-900'}`}>
                  {result.is_correct ? <CheckCircle2 className="h-5 w-5 shrink-0" /> : <XCircle className="h-5 w-5 shrink-0" />}
                  <div>
                    <div className="font-bold text-sm">{result.message}</div>
                    <div className="text-xs mt-1 opacity-80">
                      正确答案：{result.city_name}{result.target_name ? ` · ${result.target_name}` : ''}
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
