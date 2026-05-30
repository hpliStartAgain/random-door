/// <reference types="vite/client" />

declare module 'pannellum-react' {
  import React from 'react';
  export const Pannellum: React.FC<{
    width?: string;
    height?: string;
    image: string;
    pitch?: number;
    yaw?: number;
    hfov?: number;
    autoLoad?: boolean;
    onLoad?: () => void;
    [key: string]: unknown;
  }>;
}
