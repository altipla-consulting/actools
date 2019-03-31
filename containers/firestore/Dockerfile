
FROM launcher.gcr.io/google/debian9:latest

RUN apt-get update && \
    apt-get -y install curl openjdk-8-jre

RUN cd /opt && \
    curl https://storage.googleapis.com/firebase-preview-drop/emulator/cloud-firestore-emulator-v1.4.2.jar --output emulator.jar

CMD ["java", "-jar", "/opt/emulator.jar", "--host", "0.0.0.0", "--port", "12000"]
