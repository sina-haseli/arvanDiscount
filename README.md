# arvanDiscount

This project is my code challenge for ArvanCloud company interview.<br>
It uses redis pub/sub as a queue for communicate with wallet to apply credit.<br>
there is 3 api for discount:<br>
1. `/api/voucher/create`: create a voucher<br>
2. `/api/voucher/redeem`: redeem a voucher<br>
3. `/api/voucher/:voucherCode/used`: check if a voucher is used and show the list of users<br>
For run project you can use the following command:<br>
`$ go run main.go serve` or `$ discount serve`