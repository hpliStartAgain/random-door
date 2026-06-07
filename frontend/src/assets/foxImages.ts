// 七张狐狸主题贴图路径映射
// 图片文件请放置在 public/fox/ 目录下，文件名与下方一致
export const foxImages = {
  fortuneTable:  '/fox/神秘的旅行与占卜桌.png',
  compass:       '/fox/可爱狐狸与神奇罗盘.png',
  magicDoor:     '/fox/可爱狐狸开启魔法门口.png',
  cityCard:      '/fox/可爱狐狸与秘密卡片.png',
  ticketReveal:  '/fox/庆祝旅行的小狐狸.png',
  passportStamp: '/fox/萌狐护照印章时光.png',
  postcardBoard: '/fox/旅行记忆拼贴板.png',
} as const;

export type FoxImageKey = keyof typeof foxImages;
