{
  WEAVE_GITOPS_FEATURE_TELEMETRY: 'true',
  adminUser: {
    create: true,
    passwordHash: 'something',
    username: 'admin',
  },
  ingress: {
    enabled: true,
    className: 'alb',
  },
}
