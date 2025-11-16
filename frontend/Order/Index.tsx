import { Modal, notification } from "antd";
import OrderNew from "./New";
import OrderShow from "./Show/index";
import SellOff from "./SellOff/Index";
import HandingOff from "./HandingOff/index";
import Printer from "./Printer/index";
import Device from "./Device/index";
import Message from "./Message/index";
import { useCallback, useEffect, useState } from "react";
import { Bridge, IMessageTypeAudio } from "@/flutter";
import {
  ActivePageContext,
  activePageContextType,
} from "./contexts/activePageContext";
import {
  ReloadInitContext,
  shouldReloadInitContextType,
} from "./contexts/reloadInitContext";
import api from "@/services/order";
import SellOffApi from "@/services/sell-off";
import { InitDataContext } from "./contexts/initDataContext";
import { PortContext, portContextType } from "./contexts/portContext";
import styles from "./index.module.scss";
import Logo from "@/assets/images/logo.jpg";
import zhuotai from "@/assets/images/zhuotai.png";
import zhuotaiActive from "@/assets/images/zhuotai-active.png";
import diancan from "@/assets/images/diancan.png";
import diancanActive from "@/assets/images/diancan-active.png";
import dingdan from "@/assets/images/dingdan.png";
import dingdanActive from "@/assets/images/dingdan-active.png";
import jiedan from "@/assets/images/jiedan.png";
import jiedanActive from "@/assets/images/jiedan-active.png";
import guqing from "@/assets/images/guqing.png";
import guqingActive from "@/assets/images/guqing-active.png";
import jiaoban from "@/assets/images/jiaoban.png";
import jiaobanActive from "@/assets/images/jiaoban-active.png";
import dayin from "@/assets/images/dayin.png";
import dayinActive from "@/assets/images/dayin-active.png";
import loginout from "@/assets/images/login-out.png";
import loginoutActive from "@/assets/images/login-out-active.png";
import { useWebSocket, WSResponse } from "./hooks/useWebSocket";
import {
  CartReducer,
  CartContext,
  CartContextType,
} from "./contexts/cartContext";
import { useImmerReducer } from "use-immer";
import { CartAction } from "./contexts/Types";

const Order = () => {
  const [initData, setInitData] = useState<API.Order.InitData>();
  const [sellOffData, setSellOffData] = useState<API.SellOff.SellOffData>(); //沽清总列表数据
  const [tabList] = useState([
    // tab列表
    // {
    //   name: "桌台",
    //   icon: zhuotai,
    //   activeIcon: zhuotaiActive,
    //   page: "zhuotai",
    // },
    {
      name: "点餐",
      icon: diancan,
      activeIcon: diancanActive,
      page: "order-new",
    },
    {
      name: "订单",
      icon: dingdan,
      activeIcon: dingdanActive,
      page: "order-show",
    },
    // {
    //   name: "桌台接单",
    //   icon: jiedan,
    //   activeIcon: jiedanActive,
    //   page: "jiedan",
    // },
    {
      name: "估清",
      icon: guqing,
      activeIcon: guqingActive,
      page: "sell-off",
    },
    {
      name: "交接班",
      icon: jiaoban,
      activeIcon: jiaobanActive,
      page: "handing-off",
    },
    // {
    //   name: "打印机管理",
    //   icon: dayin,
    //   activeIcon: dayinActive,
    //   page: "printer",
    // },
    {
      name: "设备管理",
      icon: dayin,
      activeIcon: dayinActive,
      page: "device",
    },
    {
      name: " 退出登录",
      icon: loginout,
      activeIcon: loginoutActive,
      page: "",
    },
  ]);
  const [activePage, setActivePage] = useState<string>("order-new");
  const [active, setActive] = useState<string>("all");
  const [shouldReloadInit, setshouldReloadInit] = useState<boolean>(false);
  const [portObject, setPortObject] = useState<SerialPort | undefined>();
  // 显示确认退出弹框
  const [showLoginOutModal, setShowLoginOutModal] = useState<boolean>(false);
  // 语音播报
  const AudioPayment = (msg: string) => {
    if ("speechSynthesis" in window) {
      // HTML5中新增的API,用于将指定文字合成为对应的语音.也包含一些配置项,指定如何去阅读(语言,音量,音调)等
      const audioMessage = new SpeechSynthesisUtterance();
      audioMessage.lang = "zh-CN";
      audioMessage.text = msg;
      // 将文本值读出来
      window.speechSynthesis.speak(audioMessage);
    } else {
      console.log("浏览器不支持语音播报");
    }

    // flutter app 发出声音
    const bridge = new Bridge();
    bridge
      .sendWithResult({
        type: IMessageTypeAudio,
        params: {
          message: msg,
        },
      })
      .then((resp) => {
        // resp => { "hello": name }
      });
  };
  // 建立websocket连接，监听后台下单消息
  const onReceiveMessage = useCallback((data: WSResponse) => {
    // 如果小程序有下单，自动弹出
    if (data.order_id) {
      notification.info({
        message: `请注意`,
        description: data.title,
      });
      // 语音播报
      AudioPayment(data.title);
    }
  }, []);
  useWebSocket(onReceiveMessage);

  const [cart, dispatch] = useImmerReducer<API.Order.Cart, CartAction>(
    CartReducer,
    // 初始数据，购物车产品列表为空
    {
      products: [],
    }
  );

  const cartContextValue: CartContextType = {
    cart: cart,
    setCart: dispatch,
  };

  const activePageContextValue: activePageContextType = {
    activePage: activePage,
    setActivePage: setActivePage,
  };

  const ReloadInitContextValue: shouldReloadInitContextType = {
    reloadInit: shouldReloadInit,
    setShouldReloadInit: setshouldReloadInit,
  };
  const PortObjectContextValue: portContextType = {
    portObject: portObject,
    setPortObject: setPortObject,
  };
  // 初始化下单页面数据
  async function loadInitData() {
    const data: API.Order.InitData = await api.getInitData();
    const skus: { [key: API.Order.ProductID]: API.Order.Sku[] } = {};

    // 将产品的skus重新组织成 map[productID] = skus的形式
    // 方便子组件快速定位单一产品的规格数据
    // @todo 需要将sku属性同步返回给前端
    Object.keys(data.products).map((categoryID) => {
      const products = data.products[categoryID];
      products.map((product) => {
        skus[product.id] = product.skus;
      });
    });
    data.skuData = skus;

    setInitData(data);
  }

  // 初始化沽清数据
  async function loadSellOffData() {
    const data: API.SellOff.SellOffData = await SellOffApi.getInitData();
    setSellOffData(data);
  }

  useEffect(() => {
    loadSellOffData();
    loadInitData();
  }, [shouldReloadInit]);
  return (
    <>
      {/** 将初始数据向下传输 */}
      <InitDataContext.Provider value={initData}>
        <ReloadInitContext.Provider value={ReloadInitContextValue}>
          <ActivePageContext.Provider value={activePageContextValue}>
            <PortContext.Provider value={PortObjectContextValue}>
              <div className={styles.new_box}>
                <div className={styles.container}>
                  <div className={styles.sidebar}>
                    <img className={styles.logo} src={Logo} alt="" />
                    {tabList.map((item) => (
                      <div
                        key={item.name}
                        className={`${styles.tab_item} ${
                          activePage === item.page ? styles.active : ""
                        }`}
                        onClick={() => {
                          if (item.page === "") {
                            setShowLoginOutModal(true);
                          } else {
                            if (item.page === "device") {
                              setActive("all");
                            }
                            setActivePage(item.page);
                          }
                        }}
                      >
                        <img
                          className={styles.icon}
                          src={
                            activePage === item.page
                              ? item.activeIcon
                              : item.icon
                          }
                        />
                        <div>{item.name}</div>
                      </div>
                    ))}
                  </div>
                  <div className={styles.content}>
                    {/** 该区域都需要对购物车里的内容进行操作，将购物车内容保存到context，向下传出 */}
                    <CartContext.Provider value={cartContextValue}>
                      <div
                        style={{
                          display:
                            activePage === "order-new" ? "block" : "none",
                          height: "100%",
                        }}
                      >
                        <OrderNew initData={initData} />
                      </div>
                      {activePage === "order-show" && <OrderShow />}
                    </CartContext.Provider>
                    <div
                      style={{
                        display: activePage === "sell-off" ? "block" : "none",
                        height: "100%",
                      }}
                    >
                      <SellOff sellOffData={sellOffData} />
                    </div>
                    {activePage === "handing-off" && <HandingOff />}
                    {activePage === "printer" && <Printer />}
                    {activePage === "device" && (
                      <Device active={active} setActive={setActive} />
                    )}

                    {activePage === "message" && <Message />}
                  </div>
                </div>
              </div>
            </PortContext.Provider>
          </ActivePageContext.Provider>
        </ReloadInitContext.Provider>
      </InitDataContext.Provider>
      {/* 退出登录确认弹框 */}
      <Modal
        title="您确定退出登录嘛"
        open={showLoginOutModal}
        cancelText="取消"
        okText="确定"
        onOk={() => {
          localStorage.removeItem("token");
          localStorage.clear();
          window.location.href = "/";
        }}
        onCancel={() => setShowLoginOutModal(false)}
      ></Modal>
    </>
  );
};

export default Order;
