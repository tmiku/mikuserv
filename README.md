# mikuserv
A utility API for projects on the website. I generally test this one just using "go run ." and curl, but the Docker Compose file in the website repo should use the (very straightforward) Dockerfile here.

## endpoints
There's currently three endpoints: 
- **/mikuserv/ping**, a simple heartbeat
- **/mikuserv/contact**, backend for contact form. Takes in a JSON object with fields email, text, time, and password (password is a spam honeypot, not displayed on site). It returns a 202 if Sendgrid accepted the email, 403 if the honeypot was filled, or 500 otherwise.
- **/mikuserv/stravaToken**, backend for Strava auth. Simply takes in the auth code from Strava as an HTTP query parameter (key "code"), and communicates with Strava to exchange it for an authorization token. Sends back the body of the response received from Strava. **NOTE**: This exchanges access tokens and MUST NOT be used with HTTP, only HTTPS!!! In the current setup theres' just localhost communication between mikuserv and nginx with the reverse proxy, and nginx has HTTPS set up.

