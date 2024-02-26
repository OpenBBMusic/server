import { api } from '@/app/api';
import { Setting, SettingInfo } from '@/app/api/setting';
import { JsonCacheStorage } from '@/lib/cacheStorage';

export const settingCache = new JsonCacheStorage<SettingInfo>('bb-setting');

export class SettingInstance implements Setting {
  getInfo = async () => {
    const config = (await settingCache.get()) || {
      openMusicOrderOrigin: [],
      musicServices: [],
    };
    return {
      ...config,
      userMusicOrderOrigin:
        api.userMusicOrder.map((u) => {
          return {
            name: u.name,
            config: {},
          };
        }) || [],
    };
  };
  updateOpenMusicOrderOrigin: Setting['updateOpenMusicOrderOrigin'] = async (value) => {
    return settingCache.update('openMusicOrderOrigin', value);
  };
  updateUserMusicOrderOrigin: Setting['updateUserMusicOrderOrigin'] = async (value) => {
    return settingCache.update('userMusicOrderOrigin', value);
  };
}
