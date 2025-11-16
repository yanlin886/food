import { createContext, useContext } from "react";

export const InitDataContext = createContext<API.Order.InitData | undefined>(
  undefined
);

export type shouldReloadInitContextType = {
  reloadInit: boolean;
  setShouldReloadInit: (s: boolean) => void;
};

export const ReloadInitContext = createContext<shouldReloadInitContextType>({
  reloadInit: false,
  setShouldReloadInit: (s: boolean) => "",
});

export function useReloadInit() {
  return useContext(ReloadInitContext);
}

// 自定义hook，方便子组件直接获下单页面原始数据
export function useInitData() {
  return useContext(InitDataContext);
}
