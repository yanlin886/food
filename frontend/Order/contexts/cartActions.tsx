import { Draft } from "immer";
import { CartAction } from "./Types";
import { message } from "antd";
import { v4 as uuidv4 } from "uuid";

export function insertNewProduct(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (action.product === undefined) {
    // @todo 修改为Message
    message.warning("添加新菜品，缺少菜品参数");
    return;
  }
  const product = action.product;

  const cartProduct: API.Order.CartProduct = {
    id: product.id,
    title: product.title,
    uuid: uuidv4(),
    category_name: product.category_name,
    // @todo 这里是默认的sku id。需要跟后端确认下
    sku_id: product.skus[0].id,
    weight: product.weight,
    quantity: 1,
    // @todo 跟后端确认下产品单价的来源，是不是所有sku中最低的价格
    unit_price: product.amount,
    // 会员单价
    membership_unit_price: product.skus[0].membership_amount || product.amount,
    // 是否启用会员单价
    // use_membership_unit_price: product.skus[0].use_membership_amount || true,
    use_membership_unit_price: draft.useMemberPriceAll
      ? product.skus[0].use_membership_amount || true
      : false,
    // 绑定的配菜
    bind: product.bind,
    unit: product.unit,
    amount: product.amount,
  };

  // 如果提供了额外的sku，需要将sku值算进去
  if (action.selectedSku !== undefined) {
    cartProduct.sku_id = action.selectedSku.id;
    // 产品的规格改成选中的sku的名字
    cartProduct.attribute = action.selectedSku.name;
    // 将产品的价格修改为选中的sku的价格
    cartProduct.unit_price = action.selectedSku.amount;
    cartProduct.membership_unit_price =
      action.selectedSku.membership_amount || action.selectedSku.amount;
    cartProduct.use_membership_unit_price = draft.useMemberPriceAll
      ? action.selectedSku.use_membership_amount || true
      : false;
  }
  // 默认新增一个主菜
  draft.products.push(cartProduct);
  // 将活跃菜品，改为最后一个
  draft.activeProductIndex = draft.products.length - 1;
}

export function reverseCart(draft: Draft<API.Order.Cart>, action: CartAction) {
  if (action.reverseCart === undefined) {
    message.warning("反结账，参数不足");
    return;
  }
  // 将要放到购物车的菜品数据
  const cartProduct: API.Order.CartProduct[] = [];

  action.reverseCart.order_items.map((item) => {
    // 主菜数据
    const temp: API.Order.CartProduct = {
      id: item.product_id,
      title: item.product_name,
      // category_name: item.category_name,
      uuid: uuidv4(),
      sku_id: item.sku_id,
      quantity: item.quantity,
      unit_price: item.regular_amount,
      membership_unit_price: action.reverseCart?.membership_phone
        ? item.member_amount
        : item.regular_amount,
      attribute: item.sku_name,
      use_membership_unit_price: true,
    };
    // 主菜绑定的配菜数据
    action.initData?.products?.[item.category_id].map((product) => {
      if (product.id === item.product_id) {
        temp.bind = product.bind;
      }
    });
    // 主菜已选中的配菜数据
    if (item.ingredients) {
      temp.extra = item.ingredients.map((ingredient) => {
        // 找到配菜的分类
        for (const key in temp.bind?.ingredients) {
          temp.bind?.ingredients[key].map((groupItem) => {
            if (groupItem.id === ingredient.ingredient_id) {
              action.initData?.ingredient_categories.map((category) => {
                if (category.id === key) {
                  ingredient.type = category.name;
                }
              });
            }
          });
        }
        // 返回选中的配菜数据
        return {
          id: ingredient.ingredient_option_id,
          category: ingredient.ingredient_name,
          name: ingredient.ingredient_option_name,
          unit_price: ingredient.amount,
          quantity: ingredient.quantity,
          type: ingredient.type,
        };
      });
    }
    // 将主菜数据加到购物车中
    cartProduct.push(temp);
  });
  // 将购物车中的菜品数据，替换为反结账的菜品数据
  draft.products = cartProduct;
}

// 更新购物中一个菜品的选中sku
export function updateCartProductSku(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (
    action.cartProduct === undefined ||
    action.selectedSku === undefined ||
    action.cartProductActiveIndex === undefined
  ) {
    message.warning("更新购物车中菜品的规格，参数不足");
    return;
  }

  draft.products[action.cartProductActiveIndex] = {
    ...action.cartProduct,
    unit_price: action.selectedSku.amount,
    membership_unit_price: action.selectedSku.amount,
    attribute: action.selectedSku.name,
    sku_id: action.selectedSku.id,
  };
}

// 更新购物中一个菜品的数量
export function updateCartProductQuantity(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (
    action.cartProduct === undefined ||
    action.updatedQuantity === undefined ||
    action.cartProductActiveIndex === undefined
  ) {
    message.warning("更新购物车中菜品数量，参数不足");
    return;
  }

  if (action.updatedQuantity < 1) {
    // 如果数量小于1，直接删除
    draft.products.splice(action.cartProductActiveIndex, 1);

    // 将活跃菜品改为第一个
    draft.activeProductIndex = 0;
  } else {
    draft.products[action.cartProductActiveIndex] = {
      ...action.cartProduct,
      quantity: action.updatedQuantity,
    };
  }
}
// 更新购物中一个菜品的重量
export function updateCartProductWeight(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (
    action.cartProduct === undefined ||
    action.updatedWeiht === undefined ||
    action.cartProductActiveIndex === undefined
  ) {
    message.warning("更新购物车中菜品数量，参数不足");
    return;
  }

  if (action.updatedWeiht === "0") {
    // 如果重量等于0，直接删除
    draft.products.splice(action.cartProductActiveIndex, 1);

    // 将活跃菜品改为第一个
    draft.activeProductIndex = 0;
  } else {
    draft.products[action.cartProductActiveIndex] = {
      ...action.cartProduct,
      weight: action.updatedWeiht,
    };
  }
}

// 更新购物中活跃菜品，比如原来第一个菜品是选中的状态，现在选了第三个
export function updateCartActiveProduct(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (action.newProductIndex === undefined) {
    message.warning("更新购物车选中的菜品，参数不足");
    return;
  }

  // 判断下新的索引 < 购物车中产品数量
  if (
    action.newProductIndex > draft.products.length - 1 ||
    draft.products[action.newProductIndex] === undefined
  ) {
    message.warning("选中的菜品超出范围");
  }

  draft.activeProductIndex = action.newProductIndex;
}
// 更新某个产品的配料
export function updateCartProductExtra(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (!action.cartProduct || action.cartProductActiveIndex === undefined) {
    message.warning("更新配料，参数不足");
    return;
  }

  draft.products[action.cartProductActiveIndex] = {
    ...action.cartProduct,
    extra: action.newExtra,
  };
}
// 更新购物车的优惠
export function updateCartDiscount(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  draft.discount = action.cartDiscount;
}
// 清空购物车
export function clearCart(draft: Draft<API.Order.Cart>) {
  // 清空商品
  draft.products = [];
  // 清空优惠
  draft.discount = undefined;
  draft.valid_property = undefined;
  draft.note = "";
}

// 积分抵扣
export function updateCartPoint(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  draft.points_quantity = action.cartPoint?.points_quantity;
  draft.points_cash_out_amount = action.cartPoint?.points_amount;
}

// 修改是否使用会员价（针对某个菜品）
export function updateCartUseMembership(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (!action.cartProduct || action.cartProductActiveIndex === undefined) {
    message.warning("更新会员价，参数不足");
    return;
  }
  draft.products[action.cartProductActiveIndex] = {
    ...action.cartProduct,
    use_membership_unit_price: action.useMemberPrice,
  };
}

// 修改是否使用会员价（针对所有菜品）
export function updateCartUseMembershipAll(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  const temp: API.Order.CartProduct[] = [];
  draft.useMemberPriceAll = action.useMemberPriceAll;
  draft.products.map((item) => {
    temp.push({
      ...item,
      use_membership_unit_price: action.useMemberPriceAll,
    });
  });
  draft.products = temp;
}

// 更新优惠券信息
export function updateCartCouponReduce(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  draft.valid_property = action.couponReduce;
}

// 更新指定商品兑换券
export function updateCartCouponExchange(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  // 将兑换券信息加到购物车中
  draft.valid_property = action.couponReduceProduct;
  // 在所有菜品中，找到兑换券对应的菜品，将兑换券信息加到菜品中
  draft.products.map((item, index) => {
    if (item.uuid === action.couponReduceProduct?.exchange_product?.uuid) {
      // 如果数量大于1，需要拆分
      if (item.quantity > 1 && action.couponReduceProduct) {
        // 1. 当前菜品数量减1
        item.quantity = item.quantity - 1;
        // 2. 构造一个新的菜品，数量为1，去掉会员价并将兑换券信息加到新菜品中
        const newItem = { ...item };
        newItem.quantity = 1;
        newItem.use_membership_unit_price = false;
        newItem.coupon_exchange = action.couponReduceProduct;
        // 3. 将新的菜品插入到购物车中
        draft.products.splice(index + 1, 0, newItem);
        // 4. 将活跃菜品，改为最后一个
        draft.activeProductIndex = draft.products.length - 1;
        return;
      }
      // 数量为1， 只需将该菜品的会员价删除，并将兑换券信息加到菜品中
      item.use_membership_unit_price = false;
      item.coupon_exchange = action.couponReduceProduct;
    }
  });
}

// 删除指定商品兑换券
export function deleteCartCouponExchange(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  draft.valid_property = undefined;
  draft.products.map((item) => {
    if (item.uuid === action.couponReduceProduct?.exchange_product?.uuid) {
      item.use_membership_unit_price = false;
      item.coupon_exchange = undefined;
    }
  });
}

// 挂单取单，覆盖当前购物车
export function overlayCart(draft: Draft<API.Order.Cart>, action: CartAction) {
  if (action.holdingCart) {
    // 直接覆盖当前购物车

    // 挂单只针对没有登陆的会员。如果登陆了会员，则无法使用挂单
    // 所以这里仅需要考虑，产品+整单优惠信息即可。无需考虑优惠券，会员价等其他购物车中
    // 可能包含的营销数据
    draft.products = action.holdingCart?.cart?.products;
    draft.discount = action.holdingCart?.cart?.discount;
  }
}

// 整单备注
export function updateOrderNote(
  draft: Draft<API.Order.Cart>,
  action: CartAction
) {
  if (action.order_note) {
    draft.note = action.order_note;
  }
}
