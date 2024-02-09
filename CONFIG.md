# Environment Variables

## Version

Version info

 - `SOURCE_BRANCH` - 
 - `SOURCE_COMMIT` - 
 - `IMAGE_NAME` - 

## Config

 - `ENV` (default: `production`) - Env must be local, development, test or production
 - `TIMEZONE` (default: `Europe/Brussels`) - 
 - `HOST` - 
 - `PORT` (default: `3000`) - 
 - `ADMINS` (comma-separated, **required**, non-empty) - 
 - 
   - `STORAGE_BACKEND` (default: `s3`) - 
   - `STORAGE_CONN` (**required**, non-empty) - 
 - 
   - `REPO_CONN` (**required**, non-empty) - 
 - 
   - `OIDC_URL` (**required**, non-empty) - 
   - `OIDC_ID` (**required**, non-empty) - 
   - `OIDC_SECRET` (**required**, non-empty) - 
   - `OIDC_REDIRECT_URL` (**required**, non-empty) - 
 - 
   - `COOKIE_SECRET` (**required**, non-empty) - 
 - `MAX_FILE_SIZE` (default: `2000000000`) - 
