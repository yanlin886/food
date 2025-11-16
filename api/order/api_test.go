package order

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"jihulab.com/yanlin/food-api/pkg/auth"
	"jihulab.com/yanlin/food-api/pkg/cache"
	"jihulab.com/yanlin/food-api/pkg/constants"
	"jihulab.com/yanlin/food-api/pkg/core"
	"jihulab.com/yanlin/food-api/pkg/entity"
	"jihulab.com/yanlin/food-api/pkg/errors"
	"jihulab.com/yanlin/food-api/pkg/logger"
	"jihulab.com/yanlin/food-api/pkg/mq"
	"jihulab.com/yanlin/food-api/pkg/order"
	"jihulab.com/yanlin/food-api/pkg/test"
	"jihulab.com/yanlin/food-api/pkg/wechat/mocks"
	adminDependency "jihulab.com/yanlin/food-api/service.admin/dependency"
	couponDependency "jihulab.com/yanlin/food-api/service.coupon/dependency"
	membershipDependency "jihulab.com/yanlin/food-api/service.membership/dependency"
	"jihulab.com/yanlin/food-api/service.order/internal/operation"
	"jihulab.com/yanlin/food-api/service.order/internal/order/creation"
	paymentDependency "jihulab.com/yanlin/food-api/service.payment/dependency"
	pointDependency "jihulab.com/yanlin/food-api/service.point/dependency"
	productDependency "jihulab.com/yanlin/food-api/service.product/dependency"
	storeDependency "jihulab.com/yanlin/food-api/service.store/dependency"
	"jihulab.com/yanlin/food-api/service.store/external/store"
)

func (s *BaseTestSuite) TestAPI() {
	var tests = []test.APITestCase{
		// *** ROUTE: /web/order/product
		// 创建产品订单，不带配料
		{
			Name:    TestCase1,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"table_id": "A1",
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "pending_payment",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		{
			Name:    TestCase1_2,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "pending_payment",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    testQuery1_2,
		},
		{
			Name:    TestCase1_3,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "0",
				"total_amount": "0",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "0",
						"subtotal_sku_amount": "0",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusBadRequest,
			ExpectResponse: `
			{
				"code": "error",
				"message": "门店配置不允许支付金额为零"
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    testQuery1_3,
		},
		// 创建产品订单，不带配料, 检测入会门店
		{
			Name:   TestCase1_1,
			Method: "POST",
			URL:    "/web/order/product",
			Context: map[string]interface{}{
				// 下单对应的门店id
				constants.StoreID: int64(4),
			},
			Body: `
			{
				"membership_id": "1",
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "pending_payment",
					"type": "product"
				}
			}`,
			// 将 membership 1的入会门店设定为NULL
			Prequery:    queryCase1_1,
			IgnoreField: "order_id,created_at",
		},
		// 新建充值订单 - 含赠送金额，赠送积分，现金支付方式
		{
			Name:    TestCase2,
			Method:  "POST",
			URL:     "/web/order/recharge",
			Context: nil,
			Body: `
			{
				"membership_id": "5",
				"recharge_id": "2",
				"recharge_stall":     {
					"id": "51",
					"reward_amount": "41",
					"reward_points": 13,
					"recharge_amount": "31"
				},
				"payment_method": "cash"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 31,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "recharge"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 取消充值订单
		{
			Name:         TestCase2_1,
			Method:       "PUT",
			URL:          "/web/order/1670710200224452612/cancel",
			ExpectStatus: http.StatusOK,
			Body: `
			{
				"reason": "不想要了"
			}`,
			ExpectResponse: `
			{
				"data": "ok"
			}`,
		},
		// 创建产品订单，使用积分抵扣，用10积分抵扣2元
		{
			Name:    TestCase3,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "16",
				"total_amount": "18",
				"points_quantity": 10,
				"points_cash_out_amount": "2",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"membership_id": "5",
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366401",
					"actual_amount": 16,
					"created_at": "2023-07-06 17:11:00",
					"status": "pending_payment",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 创建产品订单，使用积分抵扣，用10积分抵扣2元, 并用余额支付
		{
			Name:    TestCase3_1,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "2",
				"total_amount": "4",
				"points_quantity": 10,
				"points_cash_out_amount": "2",
				"product_skus": [
					{
						"sku_id": "3",
						"quantity": "1",
						"amount": "4",
						"subtotal_sku_amount": "4",
						"original_subtotal_amount": "4"
					}
				],
				"membership_id": "5",
				"payment_method": "balance",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366401",
					"actual_amount": 2,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 余额支付测试
		{
			Name:    TestCase4,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "0.03",
				"payment_method": "balance",
				"membership_id": "5",
				"product_skus": [
					{
						"sku_id": "3",
						"quantity": "1",
						"amount": "0.03",
						"subtotal_sku_amount": "0.03",
						"original_subtotal_amount": "0.03"
					}
				],
				"total_amount": "0.03",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 0.03,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 创建产品订单，带配料
		{
			Name:    TestCase5,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "20",
				"total_amount": "20",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "20",
						"original_subtotal_amount": "20",
						"option_list": [
							{
								"option_id": "1",
								"amount": "1",
								"quantity": "2",
								"subtotal_amount": "2"
							}
						]
					}
				],
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"created_at": "2023-07-06 17:11:00",
					"actual_amount": 20,
					"status": "pending_payment",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 现金支付, 测试是否订单状态为已完成
		{
			Name:    TestCase6,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "cash",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 新建充值订单 - 含赠送金额，赠送积分，现金支付方式, 会员升级
		{
			Name:    TestCase7,
			Method:  "POST",
			URL:     "/web/order/recharge",
			Context: nil,
			// @todo 这里好像有点问题，既然前台传入总价1800，但为什么后台计算的结果是1000？
			// 两边是以后台数据为准么。如果以后端为准，前端好像不需要传入总价，以及套餐详情这些值。
			// @todo 目前测试下来，好像会员没有升级。测试case 好像是应该升级
			Body: `
			{
				"membership_id": "5",
				"recharge_id": "2",
				"recharge_stall":     {
					"ID": "51",
					"reward_amount": "41",
					"reward_points": 13,
					"recharge_amount": "31"
				},
				"payment_method": "cash"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 31,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "recharge"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    addUpgradeCaseInitSql,
		},
		//  新建产品订单 - 含不存在的代金优惠券
		{
			Name:    TestCase8,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"membership_id": "5",
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"coupon_property": {
					"membership_coupon_id": "120",
					"reduce_amount": "12.78",
					"category": "money_reduce",
					"name": "满200减50券"
				},
				"payment_method": "cash",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusBadRequest,
			ExpectResponse: `
			{
				"code": "error",
				"message": "优惠券没有找到，或者已经使用，请刷新订单重试"
			}`,
			IgnoreField: "order_id,created_at",
		},
		//  新建产品订单 - 含存在的代金优惠券
		{
			Name:    TestCase9,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			// @todo 传入的product skus 需要能够跟coupon property适配
			// 根据coupon list的test，构造下test_query里的数据。
			Body: `
			{
				"membership_id": "1",
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"coupon_property": {
					"membership_coupon_id": "1",
					"reduce_amount": "12.78",
					"category": "money_reduce",
					"name": "满200减50券"
				},
				"payment_method": "cash",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    rebuildProductCouponQuery,
		},
		//  新建产品订单 - 含存在的指定商品优惠券
		{
			Name:    TestCase10,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			// @todo 传入的product skus 需要能够跟coupon property适配
			// 根据coupon list的test，构造下test_query里的数据。
			Body: `
			{
				"membership_id": "1",
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "20",
						"uuid": "xx-yy-zz"
					}
				],
				"coupon_property": {
					"membership_coupon_id": "2",
					"reduce_amount": "12.78",
					"category": "product_exchange",
					"name": "9.9元喝咖啡",
					"exchange_product": {
						"amount": "9.9",
						"uuid": "xx-yy-zz"
					}
				},
				"payment_method": "cash",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    rebuildProductCouponQuery,
		},

		// 小程序创建产品订单，不带配料，前后端价格一致
		{
			Name:    TestCase11,
			Method:  "POST",
			URL:     "/wechat/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "pending",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "pending_payment",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 小程序创建产品订单，不带配料，前后端价格一致
		{
			Name:    TestCase11_1,
			Method:  "POST",
			URL:     "/wechat/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "balance",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    queryTestCase11_1,
		},
		// 获取产品订单 - 产品，管理端, 含优惠券，会员价
		{
			Name:         TestCase13,
			Method:       "GET",
			URL:          "/admin/order/1670710200224452608",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_info": {
						"id": "1670710200224452608",
						"status": "pending_delivery",
						"total_amount": 0.03,
						"actual_amount": 0.03,
						"type": "product",
						"note": "note1",
						"serial_number": "12",
						"category": "in_store_dining",
						"source": "web",
						"table_name": "A1",
						"created_at": "2023-06-20 00:28:46"
					},
					"order_item_info": {
						"order_items": [
							{
								"id": "1",
								"title": "羊肉面",
								"attributes": "默认",
								"quantity": 1,
								"unit_amount": 0.03,
								"total_amount": 0.03,
								"original_total_amount": 0.05,
								"reduce_amount": {
									"membership_amount": 0.18,
									"coupon": {
										"name": "3分钱喝咖啡",
										"amount": 0.02
									}
								}
							}
						]
					},
					"pay_info": {
						"id": "1",
						"method": "balance",
						"status": "paid",
						"pay_time": "2023-06-20 00:28:46",
						"amount": 0.03
					},
					"membership_info": {
						"name": "name5",
						"phone": "19944447116"
					},
					"coupon": {
						"membership_coupon_id": "1",
						"reduce_amount": "12.89",
						"category": "money_reduce",
						"name": "满50元减10元"
					},
					"membership_preferential_amount": 0.18,
					"store_info": {
						"name": "金鹰门店1"
					}
				}
			}`,
			IgnoreField: "pay_time",
		},
		// 小程序产品订单列表
		{
			Name:         TestCase14,
			Method:       "GET",
			URL:          "/wechat/order/product?current_page_index=1&page_size=10&status=all",
			Context:      nil,
			Body:         ``,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"total": 3,
					"records": [
						{
							"order_id": "1668915764205195262",
							"first_product_name": "红烧牛肉面",
							"actual_amount": 11,
							"images": ["https://www.jeck.tang.com/003.gif", "https://www.jeck.tang.com/003.gif", "https://www.jeck.tang.com/003.gif"],
							"item_count": 3,
							"created_at": "2023-09-12 01:38:19",
							"status": "completed",
							"category": "in_store_dining",
							"source": "wechat",
							"pay_method": "wechat"
						},
						{
							"order_id": "1668915764205195264",
							"first_product_name": "蛋炒饭",
							"actual_amount": 37,
							"images": ["https://www.jeck.tang.com/003.gif"],
							"item_count": 1,
							"created_at": "2023-01-23 01:38:19",
							"status": "completed",
							"category": "in_store_dining",
							"source": "web",
							"pay_method": "cash"
						},
						{
							"order_id": "1668915764205195263",
							"first_product_name": "红烧牛肉面",
							"actual_amount": 12,
							"images": ["https://www.jeck.tang.com/003.gif"],
							"item_count": 1,
							"created_at": "2022-06-15 01:38:19",
							"status": "completed",
							"category": "in_store_dining",
							"source": "wechat",
                            "pay_method": "balance"
						}
					]
				}
			}`,
			IgnoreField: "created_at",
			Prequery:    WechatOrderList,
		},
		// 获取产品订单 - 产品，管理端，会员价
		{
			Name:         TestCase15,
			Method:       "GET",
			URL:          "/admin/order/1670710200224452609",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_info": {
						"id": "1670710200224452609",
						"status": "pending_payment",
						"total_amount": 0.03,
						"actual_amount": 0.03,
						"type": "product",
						"note": "note2",
						"category": "in_store_dining",
						"source": "web",
						"serial_number": "3",
						"created_at": "2023-06-20 00:28:46"
					},
					"order_item_info": {
						"order_items": [
							{
								"id": "2",
								"title": "羊肉面",
								"attributes": "默认",
								"quantity": 1,
								"unit_amount": 0.03,
								"total_amount": 0.03,
								"original_total_amount": 0.06,
								"reduce_amount": {
									"membership_amount": 0.19
								}
							}
						]
					},
					"pay_info": {
						"id": "2",
						"method": "balance",
						"status": "pending",
						"pay_time": "2023-06-20 00:28:46",
						"amount": 0.03
					},
					"membership_info": {
						"name": "name5",
						"phone": "19944447116"
					},
					"point_cash_out_info": {
						"points_quantity": 123,
						"points_cash_out_amount": 5.45,
						"obtain_points": 12
					},
					"membership_preferential_amount": 0.19,
					"store_info": {
						"name": "金鹰门店1"
					}
				}
			}`,
			IgnoreField: "pay_time",
		},
		// 微信小程序订单详情
		{
			Name:         TestCase17,
			Method:       "GET",
			URL:          "/wechat/order/1670710200224452609",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_info": {
						"id": "1670710200224452609",
						"status": "pending_payment",
						"total_amount": 0.03,
						"actual_amount": 0.03,
						"type": "product",
						"note": "note2",
						"category": "in_store_dining",
						"source": "web",
						"serial_number": "3",
						"created_at": "2023-06-20 00:28:46"
					},
					"order_item_info": {
						"order_items": [
							{
								"id": "2",
								"title": "羊肉面",
								"attributes": "默认",
								"quantity": 1,
								"unit_amount": 0.03,
								"total_amount": 0.03,
								"original_total_amount": 0.06,
								"reduce_amount": {
									"membership_amount": 0.19
								}
							}
						]
					},
					"pay_info": {
						"id": "2",
						"method": "balance",
						"status": "pending",
						"pay_time": "2023-06-20 00:28:46",
						"amount": 0.03
					},
					"membership_info": {
						"name": "name5",
						"phone": "19944447116"
					},
					"point_cash_out_info": {
						"points_quantity": 123,
						"points_cash_out_amount": 5.45,
						"obtain_points": 12
					},
					"membership_preferential_amount": 0.19,
					"store_info": {
						"name": "金鹰门店1"
					}
				}
			}`,
			IgnoreField: "pay_time",
		},
		{
			Name:    TestCase18,
			Method:  "POST",
			URL:     "/wechat/order/recharge",
			Context: nil,
			Body: `
			{
				"recharge_id": "2",
				"recharge_stall":     {
					"ID": "51",
					"reward_amount": "41",
					"reward_points": 13,
					"recharge_amount": "31"
				},
				"payment_method": "cash"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 31,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "recharge"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    addUpgradeCaseInitSql,
		},
		{
			Name:         TestCase19,
			Method:       "GET",
			URL:          "/wechat/order/1670710200224452609",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_info": {
						"id": "1670710200224452609",
						"status": "pending_payment",
						"total_amount": 0.03,
						"actual_amount": 0.03,
						"type": "product",
						"category": "in_store_dining",
						"source": "web",
						"note": "note2",
						"created_at": "2023-06-20 00:28:46"
					},
					"order_item_info": {
						"order_items": [
							{
								"id": "2",
								"title": "羊肉面",
								"attributes": "默认",
								"quantity": 1,
								"unit_amount": 0.03,
								"total_amount": 0.03,
								"original_total_amount": 0.06,
								"reduce_amount": {
									"membership_amount": 0.19
								}
							}
						]
					},
					"pay_info": {
						"id": "2",
						"method": "balance",
						"status": "pending",
						"pay_time": "2023-06-20 00:28:46",
						"amount": 0.03
					},
					"membership_info": {
						"name": "name5",
						"phone": "19944447116"
					},
					"point_cash_out_info": {
						"points_quantity": 123,
						"points_cash_out_amount": 5.45,
						"obtain_points": 12
					},
					"membership_preferential_amount": 0.19,
					"store_info": {
						"name": "金鹰门店1"
					},
					"service_fee_info":{
						"service_fees":[
							{
								"amount":100, 
								"name":"服务费1"
							}
						]
					}
				}
			}`,
			Prequery: orderNotSerialNumberSql,
		},
		{
			Name:         TestCase20,
			Method:       "GET",
			URL:          "/web/order/1670710200224452609",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_info": {
						"id": "1670710200224452609",
						"status": "pending_payment",
						"total_amount": 0.03,
						"actual_amount": 0.03,
						"type": "product",
						"category": "in_store_dining",
						"source": "web",
						"note": "note2",
						"created_at": "2023-06-20 00:28:46"
					},
					"order_item_info": {
						"order_items": [
							{
								"id": "2",
								"title": "羊肉面",
								"attributes": "默认",
								"quantity": 1,
								"unit_amount": 0.03,
								"total_amount": 0.03,
								"original_total_amount": 0.06,
								"reduce_amount": {
									"membership_amount": 0.19
								}
							}
						]
					},
					"pay_info": {
						"id": "2",
						"method": "balance",
						"status": "pending",
						"pay_time": "2023-06-20 00:28:46",
						"amount": 0.03
					},
					"membership_info": {
						"name": "name5",
						"phone": "19944447116"
					},
					"point_cash_out_info": {
						"points_quantity": 123,
						"points_cash_out_amount": 5.45,
						"obtain_points": 12
					},
					"membership_preferential_amount": 0.19,
					"store_info": {
						"name": "金鹰门店1"
					},
					"service_fee_info":{
						"service_fees":[
							{
								"amount":100, 
								"name":"服务费1"
							}
						]
					}
				}
			}`,
			Prequery: orderNotSerialNumberSql,
		},
		{
			Name:         TestCase21,
			Method:       "POST",
			URL:          "/web/order/product",
			ExpectStatus: http.StatusOK,
			Body: `
			{
				"actual_amount": "12",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_total_amount": "24",
						"membership_reduce_amount": "6"
					}
				],
				"total_amount": "18",
				"payment_method": "cash",
				"membership_id": "5",
				"membership_reduce_amount": "6",
				"category": "in_store_dining"
			}`,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1685859511140618240",
					"actual_amount": 12,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    `truncate orders cascade;`,
		},
		// {
		// 	Name:         TestCase28,
		// 	Method:       "POST",
		// 	URL:          "/web/order/product",
		// 	ExpectStatus: http.StatusOK,
		// 	Body: `
		// 	{
		// 		"table_id": "1",
		// 		"service_fees": [
		// 			{
		// 				"id": "1",
		// 				"amount": "100.1"
		// 			}
		// 		],
		// 		"actual_amount": "12",
		// 		"product_skus": [
		// 			{
		// 				"sku_id": "2",
		// 				"quantity": "1",
		// 				"amount": "18",
		// 				"subtotal_sku_amount": "18",
		// 				"original_total_amount": "24",
		// 				"membership_reduce_amount": "6"
		// 			}
		// 		],
		// 		"total_amount": "18",
		// 		"payment_method": "cash",
		// 		"membership_id": "5",
		// 		"membership_reduce_amount": "6",
		// 		"category": "in_store_dining"
		// 	}`,
		// 	ExpectResponse: `
		// 	{
		// 		"data": {
		// 			"order_id": "1685859511140618240",
		// 			"actual_amount": 12,
		// 			"created_at": "2023-07-06 17:11:00",
		// 			"status": "completed",
		// 			"type": "product"
		// 		}
		// 	}`,
		// 	IgnoreField: "order_id,created_at",
		// 	Prequery:    webOrderTest,
		// },
		{
			Name:         TestCase22,
			Method:       "GET",
			URL:          "/wechat/order/balance?current_page_index=1&page_size=10",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"total": 2,
					"records": [
						{
							"description": "在线充值",
							"amount": 41,
							"created_at": "2023-07-13 01:38:19",
							"order_id": "1668915764205195265"
						},
						{
							"description": "订单支出",
							"amount": -37,
							"created_at": "2023-01-23 01:38:19",
							"order_id": "1668915764205195264"
						}
					]
				}
			}`,
			Prequery: wechatBalanceListSql,
		},
		{
			Name:         TestCase24,
			Method:       "GET",
			URL:          "/web/order/1670710200224452608/copy",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_note": "note1",
					"membership_id": "5",
					"membership_name": "name5",
					"membership_phone": "19944447116",
					"order_type":"product",
					"order_items": [
						{
							"category_id": "1",
							"product_id": "1",
							"product_name": "红烧牛肉面",
							"sku_id": "2",
							"sku_name": "中份",
							"quantity": 1,
							"member_amount": 1.23,
							"regular_amount": 1.24,
							"ingredients": [{
								"ingredient_id": "1",
								"ingredient_name": "荤菜",
								"quantity": 1,
								"amount": 3.21,
								"ingredient_option_id": "1",
								"ingredient_option_name": "鸡蛋"
							},{
								"ingredient_id": "1",
								"ingredient_name": "荤菜",
								"quantity": 3,
								"amount": 3.22,
								"ingredient_option_id": "2",
								"ingredient_option_name": "牛肉"
							}]
						}
					]
				}
			}`,
			Prequery: orderCopyTestSQL,
		},
		{
			Name:         TestCase26,
			Method:       "GET",
			URL:          "/web/order/1670710200224452609/copy",
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_note": "note1",
					"membership_id": "5",
					"membership_name": "name5",
					"membership_phone": "19944447116",
					"order_type":"product",
					"order_items": [
						{
							"category_id": "1",
							"product_id": "1",
							"product_name": "红烧牛肉面",
							"sku_id": "2",
							"sku_name": "中份",
							"quantity": 1,
							"member_amount": 1.23,
							"regular_amount": 1.24,
							"ingredients": []
						}
					]
				}
			}`,
			Prequery: orderCopyTestSQL,
		},
		// 小程序新建产品订单 - 自提订单
		{
			Name:    TestCase25,
			Method:  "POST",
			URL:     "/wechat/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_total_amount": "18"
					}
				],
				"payment_method": "cash",
				"category": "pick_up",
				"dining_time": "2023-09-15 17:00:00"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "pending_produce",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 反结账
		{
			Name:         TestCase27,
			Method:       "POST",
			URL:          "/web/order/1670710200224452608/reverse",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			Body: `
			{
				"reason": "点错了"
			} `,
			ExpectResponse: `
			{
				"data": {
					"membership_id": "5",
					"membership_name": "name5",
					"membership_phone": "19944447116",
					"order_items": [
						{
							"ingredients": [
								{
									"ingredient_id": "1",
									"ingredient_name": "荤菜",
									"ingredient_option_id": "1",
									"ingredient_option_name": "鸡蛋",
									"amount": 3.21,
									"quantity": 1
								},
								{
									"ingredient_id": "1",
									"ingredient_name": "荤菜",
									"ingredient_option_id": "2",
									"ingredient_option_name": "牛肉",
									"amount": 3.22,
									"quantity": 3
								}
							],
							"category_id": "1",
							"product_id": "1",
							"product_name": "红烧牛肉面",
							"member_amount": 1.23,
							"regular_amount": 1.24,
							"quantity": 1,
							"sku_id": "2",
							"sku_name": "中份"
						}
					],
					"order_note": "note1",
					"order_type": "product"
				}
			} `,
			Prequery: orderReverseTestSQL,
		},
		{
			Name:         TestCase29,
			Method:       "GET",
			URL:          "/admin/order?current_page_index=1&page_size=10",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"total": 2,
					"total_amount": 0.06,
					"records": [
						{
							"id": "1670710200224452608",
							"category": "in_store_dining",
							"source": "web",
							"status": "pending_delivery",
							"actual_amount": 0.03,
							"store_name": "金鹰门店1",
							"events": [
								"partial_refund",
								"pay_success"
							],
							"created_at": "2023-06-20 00:28:46"
						},
						{
							"id": "1670710200224452609",
							"category": "pick_up",
							"source": "web",
							"status": "pending_delivery",
							"actual_amount": 0.03,
							"store_name": "金鹰门店1",
							"events": [
								"pay_success",
								"reverse"
							],
							"created_at": "2023-06-20 00:28:46"
						}
					]
				}
			}`,
			Prequery: orderListTestSQL,
		},
		{
			Name:         TestCase30,
			Method:       "GET",
			URL:          "/admin/order/recharge?store_id=1&current_page_index=1&page_size=10",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"total": 1,
					"total_amount": 0.03,
					"records": [
						{
							"id": "1670710200224452612",
							"membership_name": "name3",
							"phone": "19944447114",
							"membership_card_name": "VIP",
							"membership_card_level_name": "VIP1",
							"membership_card_level": 0,
							"principal_amount": 0.03,
							"reward_amount": 0,
							"balance": 0,
							"operation_store": "金鹰门店1",
							"category": "wechat_recharge",
							"created_at": "2023-06-20 00:28:46"
						}
					]
				}
			}`,
			Prequery: queryTestCase30,
		},
		{
			Name:         TestCase31,
			Method:       "GET",
			URL:          "/web/order?only_completed_and_canceled=true&category=in_store_dining",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"records": [
						{
							"id": "1670710200224452610",
							"total_amount": 0.03,
							"status": "completed",
							"created_at": "2023-06-20 00:28:46",
							"items": [
								{
									"title": "羊肉面",
									"attributes": "默认",
									"quantity": 1,
									"unit": null
								}
							],
							"category": "in_store_dining",
							"source": "web"
						},
						{
							"id": "1670710200224452611",
							"total_amount": 0.09,
							"status": "completed",
							"created_at": "2023-06-20 00:28:46",
							"items": [
								{
									"title": "蛋炒饭",
									"attributes": "默认",
									"quantity": 1,
									"unit": null
								},
								{
									"title": "羊肉面",
									"attributes": "默认",
									"quantity": 1,
									"unit": null
								}
							],
							"category": "in_store_dining",
							"source": "wechat"
						}
					]
				}
			}`,
		},
		{
			Name:         TestCaseApp1,
			Method:       "GET",
			URL:          "/app/order?store_id=1,2,3&current_page_index=1&page_size=10&status=pending_delivery,completed",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"total": 2,
					"total_amount": 0.06,
					"records": [
						{
							"id": "1670710200224452608",
							"category": "in_store_dining",
							"source": "web",
							"status": "pending_delivery",
							"actual_amount": 0.03,
							"store_name": "金鹰门店1",
							"events": [
								"partial_refund",
								"pay_success"
							],
							"created_at": "2023-06-20 00:28:46"
						},
						{
							"id": "1670710200224452609",
							"category": "pick_up",
							"source": "web",
							"status": "pending_delivery",
							"actual_amount": 0.03,
							"store_name": "金鹰门店1",
							"events": [
								"pay_success",
								"reverse"
							],
							"created_at": "2023-06-20 00:28:46"
						}
					]
				}
			}`,
			Prequery: orderListTestSQL,
		},
		{
			Name:         TestCaseApp2,
			Method:       "GET",
			URL:          "/app/order/1670710200224452608",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_info": {
						"id": "1670710200224452608",
						"status": "pending_delivery",
						"total_amount": 0.03,
						"actual_amount": 0.03,
						"type": "product",
						"note": "note1",
						"serial_number": "12",
						"source": "web",
						"table_name": "A1",
						"category": "in_store_dining",
						"created_at": "2023-06-20 00:28:46"
					},
					"order_item_info": {
						"order_items": [
							{
								"id": "1",
								"title": "羊肉面",
								"attributes": "默认",
								"quantity": 1,
								"unit_amount": 0.03,
								"total_amount": 0.03,
								"original_total_amount": 0.05,
								"reduce_amount": {
									"membership_amount": 0.18,
									"coupon": {
										"name": "3分钱喝咖啡",
										"amount": 0.02
									}
								}
							}
						]
					},
					"pay_info": {
						"id": "1",
						"method": "balance",
						"status": "paid",
						"pay_time": "2023-06-20 00:28:46",
						"amount": 0.03
					},
					"membership_info": {
						"name": "name5",
						"phone": "19944447116"
					},
					"coupon": {
						"membership_coupon_id": "1",
						"reduce_amount": "12.89",
						"category": "money_reduce",
						"name": "满50元减10元"
					},
					"membership_preferential_amount": 0.18,
					"store_info": {
						"name": "金鹰门店1"
					}
				}
			}`,
		},
		{
			Name:         TestCaseApp3,
			Method:       "GET",
			URL:          "/app/order/recharge?store_id=1&current_page_index=1&page_size=10",
			Context:      nil,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"total": 1,
					"total_amount": 0.03,
					"records": [
						{
							"id": "1670710200224452612",
							"membership_name": "name3",
							"phone": "19944447114",
							"membership_card_name": "VIP",
							"membership_card_level_name": "VIP1",
							"membership_card_level": 0,
							"principal_amount": 0.03,
							"reward_amount": 0,
							"balance": 0,
							"operation_store": "金鹰门店1",
							"category": "wechat_recharge",
							"created_at": "2023-06-20 00:28:46"
						}
					]
				}
			}`,
			Prequery: queryTestCase30,
		},
		// 新建产品订单 - 现金支付 - 校验订单消费记录
		{
			Name:    TestCase32,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "18",
				"total_amount": "18",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "18",
						"subtotal_sku_amount": "18",
						"original_subtotal_amount": "18"
					}
				],
				"payment_method": "cash",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 18,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
		},
		// 新建产品订单 - 余额支付 - 校验订单消费记录
		{
			Name:    TestCase33,
			Method:  "POST",
			URL:     "/web/order/product",
			Context: nil,
			Body: `
			{
				"actual_amount": "0.8",
				"total_amount": "0.8",
				"product_skus": [
					{
						"sku_id": "2",
						"quantity": "1",
						"amount": "0.8",
						"subtotal_sku_amount": "0.8",
						"original_subtotal_amount": "0.8"
					}
				],
				"membership_id": "1",
				"payment_method": "balance",
				"category": "in_store_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": {
					"order_id": "1670326457995366400",
					"actual_amount": 0.8,
					"created_at": "2023-07-06 17:11:00",
					"status": "completed",
					"type": "product"
				}
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    queryTestCase33,
		},
		// 新建产品订单 - 线上支付 - 校验订单消费记录
		// 更新订单状态
		{
			Name:    TestCase35,
			Method:  "PUT",
			URL:     "/web/order/1670710200224452608/switch",
			Context: nil,
			Body: `
			{
				"current_status": "pending_delivery",
				"expected_status": "pending_dining"
			}`,
			ExpectStatus: http.StatusOK,
			ExpectResponse: `
			{
				"data": "ok"
			}`,
			IgnoreField: "order_id,created_at",
			Prequery:    queryTestCase35,
		},
	}
	e := echo.New()

	logger := logger.NewMockProvider()

	// 在订单这个场景里，我们需要真实的数据库调用，测试最终的数据库的结果是否是我们需要的结果。
	orderRepo := NewRepository(s.pool)
	mqProvider := mq.NewProvider(s.pool)

	// order event
	orderEvent := NewEvent(s.pool)

	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	// cfg, err := config.Load(stage)
	// if err != nil {
	// 	panic(err)
	// }
	cache := cache.NewProvider(s.redisClient)

	// 清除其他可能提前保存在内存中的hooks
	order.TestClearOrderCancelHooks()

	// 注入所有模块的依赖
	couponDependency.RegisterService(s.pool)
	paymentDependency.RegisterService(s.pool)
	productDependency.RegisterService(s.pool)
	membershipDependency.RegisterService(s.pool, logger)
	pointDependency.RegisterService(s.pool)
	weChatAuthSDKMocks := mocks.NewProvider(s.T())
	weChatAuthSDKMocks.On("SendSubscribeMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(func(ctx context.Context, mini entity.MiniProgram, message *entity.OrderSubscribeMessage) error {
			return nil
		})
	storeDependency.RegisterService(s.pool, cache, weChatAuthSDKMocks, logger)
	adminDependency.RegisterService(s.pool)
	// 注入自己的依赖
	RegisterConsumer(orderEvent, orderRepo)
	RegisterOrderMethod(orderRepo, s.pool, orderEvent)
	creationRepo := creation.NewRepository(s.pool)

	operationRepo := operation.NewRepository(s.pool)
	operation.RegisterHook(operationRepo)
	operation.RegisterOrderMethod(operationRepo)

	creation := creation.NewCreation(s.pool, creationRepo, mqProvider, logger)

	service := NewService(creation, s.pool, orderRepo, orderEvent, logger)

	ctx, cancel := test.NewContext()
	defer cancel()

	authProvider := auth.NewTestProvider(ctx)

	// middleware
	e.Use(errors.Middleware(logger))
	webGroup := core.WebGroup(e, authProvider)
	adminGroup := core.AdminGroup(e, authProvider)
	adminStoreGroup := core.AdminStoreGroup(e, authProvider)
	wechatGroup := core.WechatGroup(e, authProvider)
	appGroup := core.AppGroup(e, authProvider)

	storeRpc := store.NewRpc(s.pool, cache)
	orderMiddleware := NewMiddleware(storeRpc)
	RegisterHandlers(wechatGroup, appGroup, webGroup, adminGroup, adminStoreGroup, authProvider, service, orderMiddleware)

	verify := &VerifyTest{pool: s.pool}

	for _, tc := range tests {
		s.T().Run(tc.Name, func(t *testing.T) {
			if tc.Prequery != "" {
				_, err := s.testPool.Pool.Exec(ctx, tc.Prequery)
				require.Nil(t, err)
			}
			actualResponse := test.Endpoint(t, e, tc)
			// 额外的验证步骤，比如判断订单相关的营销模块的记录是否成功创建了
			verify.Switch(tc, actualResponse, t)
			s.testPool.Rebuild()
		})
	}
}
