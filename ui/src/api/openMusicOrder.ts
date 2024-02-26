import { OpenMusicOrder } from '@/app/api/openMusicOpen';
import { MusicItem, MusicOrderItem } from '@/app/api/music';
import { bb_type } from '@wails/go/models';
import { proxyMusicService } from '@/utils';

export class OpenMusicOrderInstance implements OpenMusicOrder {
  useOriginGetMusicOrder = async (url: string) => {
    const res = await proxyMusicService<bb_type.MusicOrderItem[]>({
      proxy: {
        url: '/api/open-music-order',
      },
    });

    const newList: MusicOrderItem[] = [];
    res.forEach((r) => {
      const item: MusicOrderItem = {
        ...r,
        musicList: r.musicList.map((m: MusicItem) => {
          return {
            id: m.id,
            name: m.name,
            duration: m.duration,
            origin: 'bili',
          };
        }),
      };
      newList.push(item);
    });
    return newList;
  };
}
