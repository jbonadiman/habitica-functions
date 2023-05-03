# Habitica Functions

This is a Vercel serverless function repo created specially for the pomodoro automation described [here](https://habitica.fandom.com/wiki/Android_Pomodoro_Integration).

It works as the following:<br>
[ClockWork Tomato](https://play.google.com/store/apps/details?id=net.phlam.android.clockworktomato) notifies [Tasker](https://tasker.joaoapps.com/) which in turn, calls this function.

All the pomodoro session state is managed in the function itself, with the help of a Redis instance (which I used Upstash, because of the integration with Vercel).
I plan on maybe adding more endpoints in the future if the needs arise.

Backlog:
* ~Create the vercel.json file to set function memory and more easy header control~
* Use a embedded Go database to remove Redis dependency (maybe Badger?)
