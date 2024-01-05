local base = import '../../base/albc/base.libsonnet';

std.mergePatch(base, {
  clusterName: 'staging',
  serviceAccount: {
    annotations: {
      'eks.amazonaws.com/role-arn': 'arn:aws:iam::XXX:role/some-albc-role',
    },
  },
  webhookTLS: {
    caCert: |||
      -----BEGIN CERTIFICATE-----
      MIIDQDCCAiigAwIBAgIRAPk7egvLNRiRBHoOcDuoMU4wDQYJKoZIhvcNAQELBQAw
      KjEoMCYGA1UEAxMfYXdzLWxvYWQtYmFsYW5jZXItY29udHJvbGxlci1jYTAeFw0y
      NDAxMDUwNjM5MDZaFw0zNDAxMDIwNjM5MDZaMCoxKDAmBgNVBAMTH2F3cy1sb2Fk
      LWJhbGFuY2VyLWNvbnRyb2xsZXItY2EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
      ggEKAoIBAQDwM2MD1msvcxmLeN3+IoM3qpjPrjaoOnpHTc0/hmZaFQUNINzm83P7
      KO/kmI9dsM4KshZKXqBlaWyZGTS47z2EnGgokSav3HAlLHejGZAAx27er+Rj+vku
      rsCaFDbjcfC4snrdmrB8zzWhGioEh82YRGOynekBzBZN4UY8YP5P/SKUyYoSzq0M
      WgCMifV9oSX40iMXNIWY04pvRkSZwY2ad68TI2WGFxx7kFROvjjxO+GOmq9A6ooG
      UP55RZfyz7ijYxEgvVE7YpLJeLggEnUPCa62H/4gOJG7tSMNEtZerMSdwiy6BrA+
      rZvdeKlnPPcUAVlK1Cf5AoIJ6vz1X8T7AgMBAAGjYTBfMA4GA1UdDwEB/wQEAwIC
      pDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB
      /zAdBgNVHQ4EFgQUkIgyEKf3hdGjZtGsWfiorVJAvLQwDQYJKoZIhvcNAQELBQAD
      ggEBAM4IX+79rmQsQgfVrL/mF4cmRlt1+MVFNysOB3UT7uNuCfcfJnJH2vGKXOv2
      i6IS5CNQtj/qEd6BckuzbNOEFyxV7gNaqkHo1GnjVLPqwPJfS+WicbiDgMYkWf1L
      G/rwYdCR2Hs1WOtIeBc1EjGScCfvmOvosZDY5CGxPqbFtYUx5QL6WRFI9OG+VhKx
      rSH9I4H13VeHeARhYZ1f1QdicXVl0FhyNJK7ynOM7iUvsBIO1ZwdoXrXVLVVhBTs
      YQgCtjHVPIs1rH9w9Xpao96+8FDN7n5c+qb6MnNkXXeiJgKm4+E0tdhSnIFqT9jA
      5OcgC96Dz3sq1zpEX3bQ6Ae+9MA=
      -----END CERTIFICATE-----
    |||,
    cert: |||
      -----BEGIN CERTIFICATE-----
      MIID3jCCAsagAwIBAgIQSJDty/3EiL3KQDWRIj71cjANBgkqhkiG9w0BAQsFADAq
      MSgwJgYDVQQDEx9hd3MtbG9hZC1iYWxhbmNlci1jb250cm9sbGVyLWNhMB4XDTI0
      MDEwNTA2MzkwNloXDTM0MDEwMjA2MzkwNlowLDEqMCgGA1UEAxMhYWxiYy1hd3Mt
      bG9hZC1iYWxhbmNlci1jb250cm9sbGVyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
      MIIBCgKCAQEAqoxEBGGr2I3brUXyKDseJgh2yBvUuOykkAu2MK6tJVKVlN7gfoY9
      CmFn/mf+eQORYpNJbO+6XyhHAAH88rKL+/z2xHy+YA3A5dA83betqaJWEJNNHjBO
      PupfAtWKKASURlGjXi+BV0bBCal7nyLlv3QnoFessJV83e2ZsbzQw8uCnHQMmh1p
      UTQHpMu3iNRjwTdWXAh1zDLd21l/y3kmKtYi36oWaVNj2SKSVVa51M+gaxMLTfPO
      LWTfxmvOv/P7tKPk0NprMxjgKw90K9rtltn+eE+C7OMpt5RV4aerVO5JxUtB7Umn
      QuI4OhrElRIsq3WdTJgm1EjkiXf/Tbbi/wIDAQABo4H9MIH6MA4GA1UdDwEB/wQE
      AwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
      ADAfBgNVHSMEGDAWgBSQiDIQp/eF0aNm0axZ+KitUkC8tDCBmQYDVR0RBIGRMIGO
      giZhd3MtbG9hZC1iYWxhbmNlci13ZWJob29rLXNlcnZpY2UuYWxiY4IqYXdzLWxv
      YWQtYmFsYW5jZXItd2ViaG9vay1zZXJ2aWNlLmFsYmMuc3Zjgjhhd3MtbG9hZC1i
      YWxhbmNlci13ZWJob29rLXNlcnZpY2UuYWxiYy5zdmMuY2x1c3Rlci5sb2NhbDAN
      BgkqhkiG9w0BAQsFAAOCAQEAf1jtQN4TtkOeDBYLcubPhy/QfAadmG1vydEWD7ym
      Axt72SOrdoAngtwWFaCpcAQiEGXvmzYwDk8ReAMkwDibP67tY8ybNCQge4tBQS/Z
      V6Y7xV94LRoCT5Uap7NDU9WaJ6P9lck9TV3E/4aKhn5RAx5/eCzZ/qQ830Wwcuw1
      1x9WvIx0dODFfnwSlzOPHojpUi5prwRt8YyskPiLdiscciCBmjabtAB6hwhNMlbl
      mfwoOinCoWs0wV1vu4+1z2UGbUquT6OzNwl8h5XEj9DWbcR97MInAfN7JdeRwKzU
      dhpqp4N/dosidzJNSaV815QhERVtz8ghLfzgF9dC23lVaA==
      -----END CERTIFICATE-----
    |||,
    key: |||
      -----BEGIN RSA PRIVATE KEY-----
      MIIEogIBAAKCAQEAqoxEBGGr2I3brUXyKDseJgh2yBvUuOykkAu2MK6tJVKVlN7g
      foY9CmFn/mf+eQORYpNJbO+6XyhHAAH88rKL+/z2xHy+YA3A5dA83betqaJWEJNN
      HjBOPupfAtWKKASURlGjXi+BV0bBCal7nyLlv3QnoFessJV83e2ZsbzQw8uCnHQM
      mh1pUTQHpMu3iNRjwTdWXAh1zDLd21l/y3kmKtYi36oWaVNj2SKSVVa51M+gaxML
      TfPOLWTfxmvOv/P7tKPk0NprMxjgKw90K9rtltn+eE+C7OMpt5RV4aerVO5JxUtB
      7UmnQuI4OhrElRIsq3WdTJgm1EjkiXf/Tbbi/wIDAQABAoIBAHmKhrqcpKwqxKBi
      laXoI581LvmDJAE53Dkvr4JYKdrMVP+IKnLg1cV9D3C1yhuR2F1o283/tlE0Ug0G
      Xa7UYYCOkYoL4Fxx0MO2uHnF+cRHhZBAeZgHEuwxeM57Qf8s7EKE0alAr6t0KAPr
      vcb2Rsc/TzIs6Uva7Ob3+i05g4nUZ4CjHS5D0vLtvMLxG76yI97aeBd1EH2JNSBX
      Kb1idR651QcD0Uw9pPklomb2PNPDOL7ksk4ittxa36c5W4Zdb9uv7Gvd/qfKW7TU
      8RNbUXjO/TRIJgDmEhHZzGoGBhXFx2RQkyacutjPaRjYrB20mAQ4/mRGSvyw8ZwJ
      ipT9djECgYEAxNGDdQEoOCPzaDMsrBWNtF8zSnu87EdEuuB/nxFjixL19dhLS7DL
      lWEVdGvzrv4/gN5YHRXLyf0Sseph8cyiV8zT+3cLIY/XHUg9dWY9PKdIuacciA2G
      Mwd641OJjoOQznj0nGwr/djyBSVbFXZxm+VZSYBSXGJb0PVqQQfObokCgYEA3dSH
      T1/91E3xsIk39hdRhhkyMAI6eA7tRycjP7gBhRidrHfCBaNXJ3O79BprFZ08Fobq
      7a1532/VUSyn+FM9heJdcmB0xgYngOhaTtcADtQJzPPa22H4E5zYOZ9uLRC08F7y
      hi2fCETKawVZCwgxtwtxRCzN0QzK6RJXCN0to0cCgYAmXAM3+arCDlexVk/9lhHR
      NsDDYox2rIk7tueItBXnlCF18dry/JkhGxPYZfXPhGQSFMOtn4Lhcj6DiH/gZZa6
      cARcvV3hA6zUWzEHQY7r1Fq7PFO2PJSMO4f66Fwl94RwiWh7WCXWysKYuCghbb1E
      uhWF2smykcT9W+eClyfnqQKBgAt8EllQrfKM6oNqR0RtZqIbsdZ8dwx6MVyqsQ9+
      dk2uvZMNTDVAhKWdP1DfAUZIMrEz4PvXLGUeBBxExJl9rcS9uHrQdZs+/FKXNP25
      8d3SqoM66MzM4KwbRbKOB4U7xTJgqAu8Vux1q0kpKLgCf5hrdjzCWRGGqQayFWF3
      GgLHAoGADiwz3VdqksOw5CYd5KeuFd/HugTxzxqXCZSbVW0WoCGcIQocSVYy0oSj
      afsgNs3N0UsAwVC4ukCI80/nj9AYyTXsB470v5OxMCQbl89oDsFfGbDn0wx+ckGk
      oLrS5ZjKfEDS7fuNqcaRwHtBi64tR72dpP9DuwtZ/9wD1r3Wp20=
      -----END RSA PRIVATE KEY-----
    |||,
  },
})
