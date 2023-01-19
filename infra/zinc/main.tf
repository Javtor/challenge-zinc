// Provider using local credential
provider "aws" {
  region = "us-east-2"
}
// EC2
resource "aws_instance" "ec2" {
  ami                    = "ami-00149760ce42c967b"
  instance_type          = "t2.micro"
  key_name               = aws_key_pair.default.key_name
  vpc_security_group_ids = [aws_security_group.sg.id]
  user_data_base64       = filebase64("${path.module}/user-data.sh")
}