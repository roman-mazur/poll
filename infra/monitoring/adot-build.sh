# See the docs:
# https://aws-otel.github.io/docs/setup/build-collector-as-rpm

set -e
echo "Installing AWS Distro for OpenTelemetry for $ADOT_ARCH"

(wget "https://aws-otel-collector.s3.amazonaws.com/amazon_linux/$ADOT_ARCH/v0.47.0/aws-otel-collector.rpm" && \
  wget https://aws-otel-collector.s3.amazonaws.com/aws-otel-collector.gpg && \
  rpm --import aws-otel-collector.gpg && \
  rpm --checksig aws-otel-collector.rpm && \
  rpm -Uvh ./aws-otel-collector.rpm)
