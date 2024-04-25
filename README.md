# Application Name

## Table of Contents
- [Description](#description)
- [Requirements](#requirements)
- [Requirements setup](#installation)
- [Start the app](#configuration)

## Description
This application is designed based on the assignment given to me

## Requirements
- [Golang](https://golang.org/dl/) (version 1.22.x or later)
- [SQLite](https://www.sqlite.org/download.html) installed and in your system's `PATH`
- [Stripe](https://dashboard.stripe.com/register) Account(s)
- [SupaBase](https://supabase.com/) account and google authentication provider configured

## Requirements setup
### First steps
#### Clone the repository:
   ```bash
   git clone https://github.com/mbaitar/lev-assignment.git
   cd lev-assignment
   ```
#### Create Environment file:
   ```bash
    touch .env
   ```
* Copy the content of the `.env.example` to the `.env` file

### Stripe Account(s)
1. Signup for an account on [Stripe](https://dashboard.stripe.com/register) or login if you already have an account
2. Once on the dashboard toggle `Test mode` switch on top right of the dashboard
3. Next click on the `developer` button on the left of the `Test mode` switch
4. Click on the API keys tab
5. You will see 2 keys in the standard keys list (Publishable key and Secret Key)
6. Reveal the Secret key and copy it 
7. Paste the copied key in the `.env` file => `STRIPE_API_KEY={{YOUR_KEY}}`
8. Next in left sidebar click `more`
9. Find the Connect option in the menu that pops open 
10. Press the `Get Started` button and follow the setup wizard 
11. When finished go back to the `Connect` page 
12. Go to the settings of the connect integration you created, you can find it [here](https://dashboard.stripe.com/test/settings/connect/onboarding-options/oauth)
13. Click on the `OAuth` tab and enable the `OAuth for standard accounts` option 
14. On the same page copy the `Test mode client ID`
15. Paste it in the `.env` file => `STRIPE_CLIENT_ID={{YOUR_KEY}}`
16. On the left top side you will find your company name click on it 
17. Now choose a second company name and press `next`
18. A new dashboard will open up, now follow `step 3` until `step 6`
19. Paste the copied key in the `.env` file => `STRIPE_API_KEY_TEST_USER={{YOUR_KEY}}`

### Supabase
1. Signup for an account on [Supabase](https://supabase.com/dashboard/sign-up) or login if you already have an account
2. Check your email and confirm the supabase email address
3. Once you pressed the `Confirm Email Address` button you will be redirected to the supabase dashboard
4. Click on the `+ New Project` button in the middle of the page
5. Give your project a name and password and choose a region and press `Create new project`
6. On the home page you will find a `Connect` button, press it
7. Click on the tab `App Frameworks`
8. Now you will see 2 keys `NEXT_PUBLIC_SUPABASE_URL` and `NEXT_PUBLIC_SUPABASE_ANON_KEY`
9. Copy the `NEXT_PUBLIC_SUPABASE_URL`key to `.env` => `SUPABASE_URL="{{YOUR_KEY}}"`
10. Copy the `NEXT_PUBLIC_SUPABASE_ANON_KEY`key to `.env` => `SUPABASE_SECRET="{{YOUR_KEY}}"`
11. On the left side menu click `Authentication`
12. Next under configuration on the left side of your screen you will find `Providers`
13. Once on the provider page, locate the `Email` provider and click on it
14. Disable the `Confirm email` option and save.


### Sign-in with google
1. Follow [this guide](https://supabase.com/docs/guides/auth/social-login/auth-google) on how to set login with Google

## Start the app
1. install dependencies
    ```bash
    make install
   ```
2. generate dummy data  
This will generate 1000 customers and subscriptions on your test stripe account.  
You can always stop the process if you want
    ```bash
    make data
   ```
3. Setup the database
    ```bash
    make up
   ```

4. Start the application
    ```bash
    make run
   ```
To test the buyer and seller views, you will need to create 2 accounts