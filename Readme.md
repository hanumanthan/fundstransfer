## Assumptions

1. Transaction not taken care of
2. Database recreated by app. Volumes not retained
3. All run on main go routine
4. Error validation - errors are not wrapped as user readable. thrown as what library throws
5. no unit test. integration tests

### TODO
1. Add random delay
2. Killer feature