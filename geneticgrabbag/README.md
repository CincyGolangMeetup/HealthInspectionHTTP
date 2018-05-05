# geneticgrabbag

Here's my solution for the May Cincy Go Meetup.

Alex ("searsaw") and I worked together during the meetup. Then I went home and
built a solution based on our work together.

I've tried to demonstrate some of these golang techniques:

- Placing the models in the root folder (e.g. `inspection`)
- Defining a simple interface (`InspectionRepository`) with implementations under their own directory (e.g `cincyfsp`)
- Employing the *functional options* pattern
- Developing a simple unit test that does not have access to the internals of
  the class under test
- Creating a binary (`inspect`) in a directory under `./cmd`
- Using a mock client transporter to return fake data while acting like
  calling a real REST service over the network

I owe a lot to the work of Dave Cheney and Ben Johnson.

Unfortunately, I will be unable to attend the June meetup (so bummed!) but I
can't wait for July.

Keep up the great work!

- Evan