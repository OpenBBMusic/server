import { BBMusicApp } from '@/app';
import { PcContainer } from './app/modules/container';
import { apiInstance } from './api';
import { useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';
import './style.scss';

const container = document.getElementById('root');
const root = createRoot(container!);

function Root() {
  const [initLoading, setInitLoading] = useState(false);

  const init = async () => {
    setInitLoading(true);
    try {
      await Promise.all(apiInstance.musicServices.map((s) => s.hooks?.init?.()));
    } catch (error) {
      console.log('error: ', error);
    }
    setInitLoading(false);
  };

  useEffect(() => {
    init();
  }, []);

  return (
    <div
      style={{
        width: '100%',
        height: '100%',
        borderRadius: 4,
        overflow: 'hidden',
      }}
    >
      {!initLoading && (
        <BBMusicApp apiInstance={apiInstance}>
          <PcContainer />
        </BBMusicApp>
      )}
    </div>
  );
}

root.render(<Root />);
