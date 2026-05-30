import React, { useEffect } from 'react';
import { useUserStore } from './store/useUserStore';
import { useViewStore } from './store/useViewStore';
import { MapCanvas } from './components/MapCanvas';
import { StreetViewCanvas } from './components/StreetViewCanvas';
import { Navbar } from './components/layout/Navbar';
import { Sidebar } from './components/layout/Sidebar';
import { RightDrawer } from './components/RightDrawer';

function App() {
  const initUser = useUserStore((state) => state.initUser);
  const { canvasMode } = useViewStore();

  useEffect(() => {
    initUser();
  }, [initUser]);

  return (
    <div className="relative w-screen h-screen overflow-hidden bg-background">
      <Navbar />
      
      {/* 侧边栏和底图/街景容器 */}
      <div className="absolute top-0 left-0 w-full h-full pt-[60px] flex">
        {/* 左侧固定侧栏 */}
        <Sidebar />
        
        {/* 右侧主画布区 (地图或全景) */}
        <main className="flex-1 relative bg-background overflow-hidden flex">
          {/* 地图始终在底层挂载，防止重复加载报错且能保留视角 */}
          <div className="absolute inset-0 z-0">
            <MapCanvas />
          </div>
          
          {/* 街景作为浮层覆盖在地图之上 */}
          {canvasMode === 'street' && (
            <div className="absolute inset-0 z-10 animate-in fade-in duration-500">
              <StreetViewCanvas />
            </div>
          )}
        </main>
      </div>

      {/* 全局右侧抽屉 (聊天室/图集) */}
      <RightDrawer />
    </div>
  );
}

export default App;
