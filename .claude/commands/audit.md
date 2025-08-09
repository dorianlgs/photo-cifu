Your goal is to update any vulnerable dependencies.

Only change directory in the first command

In the tests do not use yarn, only use npm

Do the following:

1. Run `cd ui && npm audit` to find vulnerable installed packages in this project
2. Run `npm audit fix` to apply updates
3. Run `npm run test_run` and verify the updates didn't break anything