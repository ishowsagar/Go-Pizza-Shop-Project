<!-- @ Workflow of our API -->

# DBModel stores --> OrderModel type -> which has Db type of \*gorm.Db

# func connectionToDB(connStr) --> creates instance of DBModel and feeds in db

# This also makes O.M in action as that type have methods on it e.g createorder

# If we pass connStr to func --> gives instance of DB.M --> Can use every model that could hace methods on it. This model acts like a source of energy point for our app

# Controller model stores db Model --> instance is created by fnc --> which needs instance of DbModel to create this instance which stores ordermodel and supplying db conn

<!-- @ next step -->

start from - https://youtu.be/8XRTAPWMO2E?list=PLra0aL87dmhcGoEaziwOjBI9KYZVKmXDc&t=6464
