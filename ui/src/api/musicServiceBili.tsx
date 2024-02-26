import {
  MusicService,
  MusicServiceAction,
  MusicServiceHooks,
  SearchItem,
  SearchParams,
  SearchType,
} from '@/app/api/musicService';
import { html2text, proxyMusicService, transformImgUrl } from '@/utils';
import { app_bili, bb_type } from '@wails/go/models';
import { settingCache } from './setting';
import { MusicItem } from '@/app/api/music';
import { SettingItem } from '@/app/modules/setting';
import { Input } from '@/app/components/ui/input';
import { useEffect, useState } from 'react';
import { Button } from '@/app/components/ui/button';
import { Switch } from '@/app/components/ui/switch';
import { message } from '@/app/components/ui/message';

class BiliMusicServiceConfigValue {
  enabled = true;
  proxyEnabled = false;
  proxyAddress = '';
  proxyToken = '';
}

const NAME = 'bili';
const CNAME = '哔哩哔哩';

class BiliAction implements MusicServiceAction {
  searchList = async (params: SearchParams) => {
    const page = params.current || 1;
    const query = {
      page: page + '',
      keyword: params.keyword,
    };

    const res = await proxyMusicService<bb_type.SearchResponse>({
      proxy: {
        url: `/api/search/${NAME}`,
        params: query,
      },
    });

    return {
      ...res,
      data: res.data.map((item) => ({
        ...item,
        name: html2text(item.name),
        cover: transformImgUrl(item.cover),
        type: item.type as SearchType,
      })),
    };
  };
  searchItemDetail = async (item: SearchItem) => {
    const info = await proxyMusicService<bb_type.SearchItem>({
      proxy: {
        url: `/api/search/${NAME}/${item.id}`,
      },
    });
    return {
      ...info,
      type: info.type as SearchType,
    };
  };
  getMusicPlayerUrl = async (music: MusicItem) => {
    return `/api/music/file/${NAME}/${music.id}`;
  };
  download = async (music: MusicItem) => {
    const a = document.createElement('a');
    a.href = `/api/music/file/${NAME}/${music.id}`;
    a.download = `${music.name}.mp4`;
    a.target = '_blank';
    a.click();
  };
}
class BiliHooks implements MusicServiceHooks {
  init = async () => {
    // await InitConfig();
  };
}

export class BiliMusicServiceInstance implements MusicService<BiliMusicServiceConfigValue> {
  name = NAME;
  cname = CNAME;
  ConfigElement = ({ onChange }: { onChange?: (v: BiliMusicServiceConfigValue) => void }) => {
    const [config, setConfig] = useState<app_bili.Config>();
    const [data, setData] = useState<BiliMusicServiceConfigValue>(
      new BiliMusicServiceConfigValue()
    );
    const loadHandler = async () => {
      const res = await proxyMusicService<app_bili.Config>({
        proxy: {
          url: `/api/config/${NAME}`,
        },
      });
      setConfig(res);
      settingCache.get().then((res) => {
        if (res) {
          const c = res.musicServices.find((m) => m.name === NAME)?.config;
          setData(c);
        }
      });
    };
    useEffect(() => {
      loadHandler();
    }, []);
    const createProps = (key: keyof typeof data, isSwitch?: boolean) => {
      return {
        [isSwitch ? 'checked' : 'value']: data[key],
        onChange: (e: any) => {
          setData((s) => ({
            ...s,
            [key]: isSwitch ? e : e.target.value,
          }));
        },
      };
    };
    const saveHandler = async () => {
      await updateMusicServicesSetting(NAME, data);
      onChange?.(data);
      message.success('已保存');
    };
    return (
      <>
        <SettingItem label='开启/关闭'>
          <Switch {...createProps('enabled', true)} />
        </SettingItem>
        <div style={{ display: data.enabled ? 'block' : 'none' }}>
          <SettingItem label='代理开启/关闭'>
            <Switch {...createProps('proxyEnabled', true)} />
          </SettingItem>
          <div style={{ display: data.proxyEnabled ? 'block' : 'none' }}>
            <SettingItem label='代理地址'>
              <Input {...createProps('proxyAddress')} />
            </SettingItem>
            <SettingItem label='代理秘钥'>
              <Input {...createProps('proxyToken')} />
            </SettingItem>
          </div>
          <div style={{ display: !data.proxyEnabled ? 'block' : 'none' }}>
            <SettingItem label='imgKey'>
              <Input
                value={config?.sign_data.img_key}
                disabled
              />
            </SettingItem>
            <SettingItem label='subKey'>
              <Input
                value={config?.sign_data.img_key}
                disabled
              />
            </SettingItem>
            <SettingItem label='UUID_V3'>
              <Input
                value={config?.spi_data.b_3}
                disabled
              />
            </SettingItem>
            <SettingItem label='UUID_V4'>
              <Input
                value={config?.spi_data.b_4}
                disabled
              />
            </SettingItem>
          </div>
        </div>
        <div style={{ marginBottom: 15, display: 'flex', gap: '15px' }}>
          <Button onClick={saveHandler}>保存</Button>
        </div>
      </>
    );
  };
  action = new BiliAction();
  hooks = new BiliHooks();
}

async function updateMusicServicesSetting(serviceName: string, data: any) {
  const setting = await settingCache.get();
  let list = setting?.musicServices || [];
  if (list.find((n) => n.name === serviceName)) {
    list = list.map((l) => {
      if (l.name === serviceName) {
        l.config = data;
      }
      return l;
    });
  } else {
    list.push({
      name: serviceName,
      config: data,
    });
  }
  settingCache.update('musicServices', list);
}
