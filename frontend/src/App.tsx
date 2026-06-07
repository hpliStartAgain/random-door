import { useEffect, useState } from 'react';
import { useUserStore } from './store/useUserStore';
import { useViewStore } from './store/useViewStore';
import { MapCanvas } from './components/MapCanvas';
import { StreetViewCanvas } from './components/StreetViewCanvas';
import { Navbar } from './components/layout/Navbar';
import { Sidebar } from './components/layout/Sidebar';
import { RightDrawer } from './components/RightDrawer';
import { AchievementPage } from './pages/AchievementPage';
import { AssetPage } from './pages/AssetPage';
import { WelcomeOverlay } from './components/overlays/WelcomeOverlay';
import { RandomCityModal } from './components/overlays/RandomCityModal';
import { Toast } from './components/Toast';
import { AdminPage } from './pages/AdminPage';

function App() {
  const initUser = useUserStore((state) => state.initUser);
  const { canvasMode, currentView, hasEntered } = useViewStore();
  const [showAdmin, setShowAdmin] = useState(false);

  useEffect(() => {
    initUser();
  }, [initUser]);

  return (
    <div className="relative w-screen h-screen overflow-hidden bg-background">
      <Navbar onAdmin={() => setShowAdmin(true)} />
      
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

          {/* 任意门 — 随机城市弹窗（覆盖在地图区域上方） */}
          {currentView === 'GAME_DICE' && (
            <div className="absolute inset-0 z-20">
              <RandomCityModal />
            </div>
          )}
        </main>
      </div>

      {/* 全局右侧抽屉 (聊天室/图集) */}
      <RightDrawer />

      {/* 全局覆盖层 (成就墙等) */}
      {useViewStore(s => s.currentView) === 'ACHIEVEMENT' && (
         <div className="absolute inset-0 z-50 bg-background/95 backdrop-blur overflow-y-auto">
           <AchievementPage />
         </div>
      )}

      {useViewStore(s => s.currentView) === 'ASSETS' && (
         <div className="absolute inset-0 z-50 bg-background/95 backdrop-blur overflow-y-auto">
           <AssetPage />
         </div>
      )}

      {/* 首屏欢迎页（未进入时全屏遮罩） */}
      {!hasEntered && <WelcomeOverlay />}

      {/* 后台管理 */}
      {showAdmin && <AdminPage onClose={() => setShowAdmin(false)} />}

      {/* 全局 Toast 通知 */}
      <Toast />
    </div>
  );
}

export default App;
