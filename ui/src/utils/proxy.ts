import { settingCache } from '@/api/setting';
import { Resp } from '@/lib/request';
import axios, { AxiosRequestConfig } from 'axios';

export class ProxyConfig {
  enabled = true;
  proxyEnabled = false;
  proxyAddress = '';
  proxyToken = '';
}

interface ProxyServiceOptions<T> {
  proxy: AxiosRequestConfig;
}

export async function proxyMusicService<T = any>({ proxy }: ProxyServiceOptions<T>) {
  const res = await axios<Resp<T>>({
    ...proxy,
  }).then((res) => res.data.data);
  return res;
}

// 是否开启了代理
export async function musicServiceEnabledProxy(origin: string) {
  const config = (await getMusicServiceConfig(origin)) as ProxyConfig;
  return config.proxyEnabled;
}

// 获取源服务的配置信息
export async function getMusicServiceConfig<T>(name: string) {
  const setting = await settingCache.get();
  const s = setting?.musicServices.find((s) => s.name === name);
  return s?.config as ProxyConfig & T;
}

export function mergeUrl(a: string, b: string) {
  let a1 = a.endsWith('/') ? a : a + '/';
  let b1 = b.startsWith('/') ? b.substring(1) : b;
  return a1 + b1;
}
