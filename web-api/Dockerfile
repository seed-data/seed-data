# Using official python runtime base image
FROM python:3.4-alpine

# Set the application directory
ADD . /app
WORKDIR /app

RUN pip install -r requirements.txt
CMD ["python", "app.py"]

# Install our requirements.txt

RUN pip install -r requirements.txt

# Copy our code from the current folder to /app inside the container
ADD . /app

# Make port 80 available for links and/or publish
EXPOSE 80

# Define our command to be run when launching the container
CMD ["gunicorn", "app:app", "-b", "0.0.0.0:80", "--log-file", "-", "--access-logfile", "-", "--workers", "4", "--keep-alive", "0"]
