import { createContext, Dispatch, useContext } from "react";
import { Draft } from "immer";
import { CartAction } from "./Types";
import {
  insertNewProduct,
  reverseCart,
  updateCartProductQuantity,
  updateCartProductWeight,
  updateCartProductSku,
  updateCartActiveProduct,
  updateCartProductExtra,
  updateCartDiscount,
  clearCart,
  updateCartPoint,
  updateCartUseMembership,
  updateCartUseMembershipAll,
  updateCartCouponReduce,
  updateCartCouponExchange,
  deleteCartCouponExchange,
  overlayCart,
  updateOrderNote,
} from "./cartActions";

// 购物车里的菜品详情作为第一级组件的state。
// 对该state的任何更新，使用reducer dispatch，而不是使用setState操作。
// 将该state跟state对应的更新dispatch使用context，向子组件传输
// https://react.dev/learn/scaling-up-with-reducer-and-context

export type CartContextType = {
  cart: API.Order.Cart | undefined;
  setCart: Dispatch<CartAction>;
};

export const CartContext = createContext<CartContextType>({
  cart: undefined,
  setCart: () => "",
});

export function useCart() {
  return useContext(CartContext);
}

export function CartReducer(draft: Draft<API.Order.Cart>, action: CartAction) {
  switch (action.type) {
    // 新增一个主菜
    case "insert_new_product":
      insertNewProduct(draft, action);
      break;
    // 反结账设置购物车
    case "reverse_cart":
      reverseCart(draft, action);
      break;
    // 更新购物车中菜品的sku
    case "update_cart_product_sku":
      updateCartProductSku(draft, action);
      break;
    // 更新选中菜品的数量
    case "update_cart_product_quantity":
      updateCartProductQuantity(draft, action);
      break;
    // 更新选中菜品的数量
    case "update_cart_product_weight":
      updateCartProductWeight(draft, action);
      break;
    // 修改活跃菜品
    case "update_active_product":
      updateCartActiveProduct(draft, action);
      break;
    // 更新选中菜品的Extra
    case "update_cart_product_extra":
      updateCartProductExtra(draft, action);
      break;
    // 更新优惠
    case "update_cart_discount":
      updateCartDiscount(draft, action);
      break;
    // 清空购物测
    case "clear_cart":
      clearCart(draft);
      break;
    // 更新积分
    case "update_cart_point":
      updateCartPoint(draft, action);
      break;
    // 更新是否使用会员价
    case "update_cart_use_member_price":
      updateCartUseMembership(draft, action);
      break;
    // 更新是否使用会员价（针对所有菜品）
    case "update_cart_use_member_price_all":
      updateCartUseMembershipAll(draft, action);
      break;
    // 更新代金券
    case "update_cart_coupon_reduce":
      updateCartCouponReduce(draft, action);
      break;
    case "update_cart_coupon_exchange":
      updateCartCouponExchange(draft, action);
      break;
    case "delete_cart_coupon_exchange":
      deleteCartCouponExchange(draft, action);
      break;
    case "overlay_cart":
      overlayCart(draft, action);
      break;
    case "update_order_note":
      updateOrderNote(draft, action);
  }
}
