// input
db.getCollection('transactions').find({"vin": {"$elemMatch": {"txid": "aa3f63a1df52b793d67788bc8193d85f06d3d1745b4851dacb2fa279b8bb28e5"}}});

// tx
db.getCollection('transactions').find({"txid": "755ee7b7c8f778c3f0907783254586a6c24dd0699f32d17a4b6a26843c7f2500"});

// output hot
db.getCollection('transactions').find({"vout": {"$elemMatch": {"scriptPubKey.addresses": {"$elemMatch": {"$eq": "n1xpEB26rSi1XkLMQFkgZ44yjSXgUpUAws"}}}}});
// output cold
db.getCollection('transactions').find({"vout": {"$elemMatch": {"scriptPubKey.addresses": {"$elemMatch": {"$eq": "2NEhic4wTnBittzJru5r6SWP8LNjHjdE7nZ"}}}}});
