# Groupie Trackers

Groupie Trackers is a web application that allows users to explore and visualize information about various bands and artists. The project uses data from a given API that consists of four key parts: artists, locations, dates, and relations. The goal is to create a user-friendly website that displays this information through various data visualizations, such as blocks, cards, tables, lists, pages, and graphics.

## Features

- Display information about bands and artists, including their names, images, start years, first album dates, and members.
- Show concert locations and dates, including past and upcoming events.
- Create interactive visualizations to link artists with their concert locations and dates.
- Implement client-server communication to trigger events/actions and retrieve data dynamically.

## API Structure

The API consists of four parts:

1. **Artists**: Contains information about bands and artists, including:
   - Name
   - Image
   - Year they began activity
   - Date of their first album
   - Members of the band

2. **Locations**: Contains data about the locations of past and upcoming concerts.

3. **Dates**: Contains data about the dates of past and upcoming concerts.

4. **Relation**: Links artists with their concert dates and locations.

## Objectives

- Build a responsive and interactive website to display artist information.
- Use data visualizations to present information clearly and engagingly.
- Implement events/actions that trigger server communication, demonstrating client-server interaction.
- Ensure the website and server are robust, handling errors gracefully without crashing.
- Write clean, maintainable code following good practices, including unit testing.

## Technologies Used

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **Data Format**: JSON

## Requirements

- The backend must be implemented in Go.
- The application should be robust, with no crashes or unhandled errors.
- Follow good coding practices, including code organization and error handling.
- Include unit tests for critical parts of the codebase.
- Use only standard Go packages.

## Setup and Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/groupie-trackers.git
   cd groupie-trackers
