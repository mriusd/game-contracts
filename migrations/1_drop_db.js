const MongoClient = require('mongodb').MongoClient;

module.exports = async function(deployer, network, accounts) {
  // Replace the following values with your MongoDB connection details
  const uri = "mongodb://localhost:27017";
  const dbName = "game";

  // Connect to MongoDB
  const client = new MongoClient(uri, { useNewUrlParser: true, useUnifiedTopology: true });
  await client.connect();

  // Drop the database
  try {
    await client.db(dbName).dropDatabase();
    console.log(`Database '${dbName}' dropped successfully.`);
  } catch (err) {
    console.log(`Error dropping database '${dbName}':`, err.message);
  } finally {
    // Close the connection
    await client.close();
  }
};
