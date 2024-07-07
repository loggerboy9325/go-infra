package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCZ1Ft5pTXMZgsjcHk6OYmYnXUM1LdKgCwxbHOGef61YiySOtUE4vY/yaGL9ca/HrVuMxEnWYTKPlnO7dkY7xZf1wscg5+oMt7fj1TvLFvSypmd1Cp/pO7heAQRJkK6KOlNocjHZ1iSPV5nvzXzu8y8i1Fi1gDIUB6lUWQCbN4hoJG3qScDoVyfdCHIvroAMNyaXWTDtaCPGKl6TjMxjS3Yqig3MxkwKdh54v9LXeVLPgAljxTTEiyFqLuM51jNGNaupcvq8cLUGQglUjJmVWNntLs6FVg78F0I8gAHlkGysFlpiOJ8fcZegfV3LgNaS0GDmj4RW8FTWtYSW7P9Dkykq7j0FtL0pOV71wfyxfcP6QrTsIh7gNnbumId/0zaV8aVb78HEKozbjL8eH626oKyu5HwLj/Va9yUrPVnUlhC1vx4S+CGi5EYaW+llB21vlARP7BoJDKYXrahH4EI3NWfmx7qVXNqJ/kaiUR/Kyj/I/nQrmbFTPFgnWwMIWoDNbM= george@george-IdeaPad-Flex-5-15ALC05"),
		})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-0195204d5dce06d99"),
			KeyName:             kp.KeyName,
		})
		if err != nil {
			return err
		}

		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("PublicIp", jenkinsServer.PublicIp)
		ctx.Export("PublicDns", jenkinsServer.PublicDns)

		return nil
	})
}
