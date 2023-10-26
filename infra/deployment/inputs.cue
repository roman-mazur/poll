package deployment

inputs: {
	tls: {
		certificate: string | *"" @tag(cert)
		private_key: string | *"" @tag(pkey)
	}
	admin: {
		secret: string | *"" @tag(admin)
	}
}
