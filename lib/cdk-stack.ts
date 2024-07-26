import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as path from "path";
// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const goLambda = new lambda.Function(this, "Schedule-Parser", {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: "bootstrap", // handler is the name of the executable
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../build/deployment.zip")
      ),
      environment: {
        OPENAI_API_KEY: process.env.OPENAI_API_KEY || "",
      },
      timeout: cdk.Duration.seconds(90),
    });

    // Define the API Gateway REST API
    const api = new apigateway.RestApi(this, "Schedule-Parser-API", {
      restApiName: "Schedule Parser Service",
      description: "This service serves parses schedules from images to JSON.",
      defaultCorsPreflightOptions: {
        allowOrigins: apigateway.Cors.ALL_ORIGINS,
      },
    });

    // Create an API key
    const apiKey = api.addApiKey("schedule-parser-user-api-key", {
      apiKeyName: "ScheduleParserUserAPIKey",
      description: "Schedule parser API key for normal usage",
    });

    // Create a usage plan and associate it with the API key
    const plan = api.addUsagePlan("UsagePlan", {
      name: "BasicUsagePlan",
      description: "A usage plan for my API",
      throttle: {
        rateLimit: 2,
        burstLimit: 3,
      },
      quota: {
        limit: 50,
        period: apigateway.Period.DAY,
      },
      apiStages: [{
        stage: api.deploymentStage,
      }],
    });

    plan.addApiKey(apiKey);

    // Integrate the Lambda function with the API Gateway
    const getLambdaIntegration = new apigateway.LambdaIntegration(goLambda, {
      requestTemplates: { "application/json": '{"statusCode": 200}' }
    });

    // Define the API endpoint and method
    const item = api.root.addResource("parse");
    item.addMethod("POST", getLambdaIntegration, {apiKeyRequired: true}); // POST /item
  }
}
