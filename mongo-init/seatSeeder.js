db = db.getSiblingDB("cinema")

db.movies.deleteMany({})
db.showtimes.deleteMany({})
db.seats.deleteMany({})
db.bookings.deleteMany({})
db.audit_logs.deleteMany({})

// Movies
const movie1 = {
  _id: "m1",
  title: "Avengers"
}

const movie2 = {
  _id: "m2",
  title: "Interstellar"
}

db.movies.insertMany([movie1, movie2])

// Showtimes
const now = new Date()

const showtime1 = {
  _id: "1",
  movie_id: movie1._id,
  movie_title: movie1.title,
  start_time: new Date(now.getTime() + 60 * 60 * 1000),
  theater: "Theater 1"
}

const showtime2 = {
  _id: "2",
  movie_id: movie2._id,
  movie_title: movie2.title,
  start_time: new Date(now.getTime() + 3 * 60 * 60 * 1000),
  theater: "Theater 2"
}


db.showtimes.insertMany([showtime1, showtime2])

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

db.seats.insertMany([
  ...generateSeats(showtime1._id),
  ...generateSeats(showtime2._id)
])

db.seats.createIndex(
  { showtime_id: 1, seat_number: 1 },
  { unique: true }
)

print("Seat seeding completed")