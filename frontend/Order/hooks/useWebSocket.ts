import { useEffect } from "react";

const token = localStorage.getItem("bearer_token") ?? "";

export type WSResponse = {
  order_id: string;
  title: string;
};

export function useWebSocket(onReceiveMessage?: (data: WSResponse) => void) {
  // https://react.dev/reference/react/experimental_useEffectEvent
  // const onMessage = useEffectEvent(onReceiveMessage);

  // 用来建立长链接，接收后端对小程序下单的
  useEffect(() => {
    let ignore = false;
    let timerID: NodeJS.Timer;
    //建立长连接
    const ws = new WebSocket(
      `${process.env.REACT_APP_WS_HOST}/web/ws/store/message/${token}`
    );
    ws.onopen = () => {
      // 定时心跳
      timerID = setInterval(() => {
        ws.send("ping");
      }, 5000);
    };

    // 读取ws消息
    ws.onmessage = (event) => {
      // 如果是心跳消息，直接忽视。
      if (event.data === "pong") {
        return;
      }
      // 读取后端发送的消息
      const data = JSON.parse(event.data);
      if (!ignore && onReceiveMessage) {
        onReceiveMessage(data);
      }
    };

    ws.onclose = () => {
      // 如果ws有心跳记录，关闭心跳。
      if (timerID) {
        clearInterval(timerID);
      }
    };

    // 关闭组件时，清除记录
    return () => {
      ignore = true;

      // 如果ws还是开启状态，关闭连接。
      if (
        ws.readyState === WebSocket.OPEN ||
        ws.readyState === WebSocket.CONNECTING
      ) {
        ws.close();
      }
    };
  }, [onReceiveMessage]);
}
