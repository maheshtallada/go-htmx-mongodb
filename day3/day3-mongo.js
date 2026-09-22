// C2
// 2.1
db.orders.find(
    {total : {$gte: 150}},
    {customer: 1, total: 1, status:1, items:1}
)

// 2.2
db.orders.find(
    {status : "paid", items: {$gte: 3}},
    {customer: 1, total: 1, status:1, items:1}
)

// 2.3
db.orders.find({
    $or: [
        { status: "paid" },
        { total: { $lt: 100 } }
    ]
})

// show dbs
// show collections
// use day3


// C3 -- copied
db.orders.find(
        {},
        { _id: 0, customer: 1, total: 1 }
        ).sort(
            { total: -1 }
        ).limit(2)