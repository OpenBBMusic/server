import { isJson } from '@/utils';
import { StateStorage } from 'zustand/middleware';
import { get, set, del } from 'idb-keyval';

/** 缓存接口 */
export const cacheStorage: StateStorage = {
  getItem: async (name: string): Promise<string | null> => {
    return (await get(name)) || '';
  },
  setItem: async (name: string, value: string): Promise<void> => {
    return await set(name, value);
  },
  removeItem: async (name: string): Promise<void> => {
    return await del(name);
  },
};

export class JsonCacheStorage<T> {
  constructor(private key: string) {}
  get = async (): Promise<T | undefined> => {
    const res = await get(this.key);
    return isJson<T>(res);
  };
  set = async (value: T): Promise<void> => {
    return await set(this.key, JSON.stringify(value));
  };
  removeItem = async (): Promise<void> => {
    return await del(this.key);
  };
  async update(key: keyof T, value: T[keyof T]): Promise<void> {
    const res = await this.get();
    const data = { ...res!, [key]: value };
    await this.set(data);
  }
}
