db = db.getSiblingDB("cinema")

db.movies.deleteMany({})
db.showtimes.deleteMany({})
db.seats.deleteMany({})

// Movies
const movie1 = {
  _id: ObjectId(),
  title: "Avengers"
}

const movie2 = {
  _id: ObjectId(),
  title: "Interstellar"
}

db.movies.insertMany([movie1, movie2])

// Showtimes
const showtime1 = {
  _id: "1",
  movie_id: movie1._id,
  start_time: new Date()
}

const showtime2 = {
  _id: "2",
  movie_id: movie2._id,
  start_time: new Date()
}

db.showtimes.insertMany([showtime1, showtime2])

// Seat layout
const rows = ["A","B","C","D","E"]

function generateSeats(showtimeId) {
  const seats = []

  rows.forEach(row => {
    for (let i = 1; i <= 5; i++) {
      seats.push({
        showtime_id: showtimeId,
        seat_number: `${row}${i}`,
        status: "AVAILABLE",
        locked_by: null,
        locked_until: null
      })
    }
  })

  return seats
}

// Insert seats
db.seats.insertMany([
  ...generateSeats(showtime1._id),
  ...generateSeats(showtime2._id)
])

// Index กันจองซ้ำ
db.seats.createIndex(
  { showtime_id: 1, seat_number: 1 },
  { unique: true }
)

print("Seat seeding completed")