# 專案說明
## 規格
設計一套搶購機制，規格：
1. 每天定點23點發送優惠券，需要用戶提前1-5分鐘前先預約，優惠券的數量為預約用戶數量的20%，
2. 搶購優惠券時，每個用戶只能搶一次，只有1分鐘的搶購時間，怎麼盡量保證用戶搶到的概率盡量貼近20%。
3. 需要考慮有300人搶和30000人搶的場景。
4. 設計表結構和索引、編寫主要程式碼。

## 設計
* 只實作預約跟搶購這兩隻 API
* 簡化設計，暫定只有這種優惠券，並且營運沒辦法自訂優惠券內容，如果要區分多種優惠券需要在 schema 上面加上 category
* 簡化設計，優惠券 開放時間 機率等 設定應該要從 DB 拿，這邊先直接 hardcode 在程式裡面
* userID 應該要從 JWT 去拿，這邊為了 demo 方便，直接寫在 header 上
* 我們在用戶預約時就決定是否要給他們優惠券，而後續的搶購階段只是根據這個結果來領取優惠券
  * 優點: 
    * 這種方法的效能相對較好，能提前計算結果，也可以確保每個用戶搶到優惠券的機率都是20%
    * 對效能影響比較大的時段能從 搶購變成預約時段，可以有效地降低對效能有較大影響的 API 的 QPS 峰值，能從 30000/60 降低成 30000/300。
    * 如果未來有更多使用者要同時參與，也能直接擴充 server/redis/mongo 機器數量來增高 QPS。
  * 缺點: 
    * 我們無法控制最終領取優惠券的用戶是誰，因此發出的優惠券數量可能會多於或少於實際領取的用戶數量的20%。
      * 缺點補救措施: 如果想要更精確的控制用戶數量，可以在預約時段結束後再多發或少發一些優惠券。


## Schema
  * 主要包含這三個欄位 couponID userID status created_at updated_at
  * 會做兩個 index，分別是 couponID 及 userID

## API 實作
* 預約
  * API - POST /coupon/reserve
  * 先檢查時間，每天 22:55 ~ 22:59 分才接受請求
  * 使用 20% 自然機率決定使用者能不能拿到優惠券，這樣優惠券的數量會接近用戶數量的 20%，每個用戶拿到的機率也會是 20%
    * (補充說明) 如果想要更精確的控制用戶數量，可以在預約時段結束後再多發或少發一些優惠券
  * 每個用戶只有一次預約的機會，透過 redis 的 NX 去避免多次重複預約，也能避免 race condition 造成一個人能重複多次抽獎
  * 最後再把優惠券訊息寫入 mongodb 中，如果寫入失敗則刪除 redis key，讓用戶有重試的機會
* 搶購
  * API - POST /coupon/snatch
  * 先檢查時間，每天 23:00 ~ 23:01 分才接受請求
  * 先去撈撈看 cache ，看看該用戶符不符合領獎資格，不符合就直接回傳
  * 如果符合領獎資格，那就把優惠券狀態改成有效的


## API文件
* 預約API
  * path: http://localhost:8080/v1/coupon/reserve
  * params:
    * header 
      * userID: {uuid}
  * respond:
    * success: bool
* 搶購API
  * path: http://localhost:8080/v1/coupon/snatch
  * params:
    * header 
      * userID: {uuid}
  * respond:
    * success: bool
    * data:
      * coupon:
        * couponID
        * userID
        * status
        * updated_at
        * created_at
      * isWinCoupon: bool


## 如何執行
`docker-compose up -d`




