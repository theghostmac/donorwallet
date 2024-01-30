# DonorWallet

Link to [Postman API Documentation](https://documenter.getpostman.com/view/19664943/2s9YysE318).

Link to [GitHub repository](https://github.com/theghostmac/donorwallet).

Link to [Render service](https://donorwallet.onrender.com).

External services worth mentioning:

- Sendgrid for "thank you" email
- Uber Zap for logging

Security measures implemented:

- GORM (against SQL injection)
- BONUS improvement: Use of JWT tokens for Authentication.

## TODO

- [x] Will be deployed on [Render](https://render.com).
- [x] Allow user create an account with basic user information
- [x] Allow a user login
- [x] Allow a user have a wallet
- [x] Allow a user create a transaction PIN
- [x] Allow a user create a donation to a fellow user (beneficiary) with adequate authorization.
- [x] Allow a user check how many donations he/she has made with adequate authorization.
- [x] Have ability to a special thank you message via email or sms or any mock communication channel, if they make two(2) or more donations.
- [x] Allow a user view all donation made in a given period of time using adequate authorization.
- [x] Allow a user view a single donation made to a fellow user (beneficiary) with adequate authorization.
- [x] Implement pagination for page and limit (donations)
