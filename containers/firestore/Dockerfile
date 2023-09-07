
FROM launcher.gcr.io/google/debian11:latest

RUN apt-get update && \
    apt-get -y install curl openjdk-11-jre

RUN cd /opt && \
    curl https://storage.googleapis.com/firebase-preview-drop/emulator/cloud-firestore-emulator-v1.4.2.jar --output emulator.jar

CMD ["java", "-jar", "/opt/emulator.jar", "--host", "0.0.0.0", "--port", "12000"]
