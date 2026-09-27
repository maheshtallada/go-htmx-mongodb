use day4

for (let i = 0; i < 1000; i++) {
    db.users.insertOne({
        name: "user" + i,
        age: 18 + (i % 50),
        city: ["NYC","LA","SF","CHI"][i % 4],
        active: i % 2 === 0
    })
}

// db.users.countDocuments()