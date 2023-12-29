local base = import '../../base/cert-manager/base.libsonnet';

std.mergePatch(base, {
  serviceAccount: {
    annotations: {
      'eks.amazonaws.com/role-arn': 'arn:aws:iam::XXX:role/some-cert-manager-role',
    },
  },
})
