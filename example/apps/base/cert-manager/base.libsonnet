{
  installCRDs: true,
  serviceAccount: {
    annotations: {
      'eks.amazonaws.com/role-arn': '',
    },
  },
  securityContext: {
    fsGroup: 1001,
  },
}
