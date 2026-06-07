/**
 * 图片分享工具：封装 canvas/blob/dataURL 转换，提供下载和复制能力。
 */

/** 将 dataURL 或 URL 下载为本地文件 */
export async function downloadImage(src: string, filename = 'capture.png'): Promise<void> {
  const a = document.createElement('a');
  if (src.startsWith('data:')) {
    a.href = src;
  } else {
    try {
      const res = await fetch(src);
      const blob = await res.blob();
      a.href = URL.createObjectURL(blob);
    } catch {
      a.href = src;
    }
  }
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

/** 将 dataURL 复制到剪贴板（不支持时返回 false，调用方可降级到 downloadImage） */
export async function copyImageToClipboard(dataUrl: string): Promise<boolean> {
  if (!navigator.clipboard || !window.ClipboardItem) return false;
  try {
    const res = await fetch(dataUrl);
    const blob = await res.blob();
    await navigator.clipboard.write([
      new ClipboardItem({ [blob.type]: blob }),
    ]);
    return true;
  } catch {
    return false;
  }
}
