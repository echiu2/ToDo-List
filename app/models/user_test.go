package models

func (ms *ModelSuite) Test_User_FullName() {
	// 1. Arrange
	tcases := []struct { // table-driven tests
		User           User
		IsErr          bool
		ExpectedResult string
	}{
		{
			User:           User{FirstName: "Joe", LastName: "Arias"},
			IsErr:          false,
			ExpectedResult: "Joe Arias",
		},
		{
			User:           User{FirstName: "Joe", LastName: "   "},
			IsErr:          true,
			ExpectedResult: "",
		},
		{
			User:           User{FirstName: "", LastName: ""},
			IsErr:          true,
			ExpectedResult: "",
		},
		{
			User:           User{FirstName: "", LastName: "Arias"},
			IsErr:          true,
			ExpectedResult: "",
		},
	}

	for index, tc := range tcases {
		// 2. Act
		fullName, err := tc.User.FullName()

		// 3. Assert
		if tc.IsErr {
			ms.Error(err, "Case#: %v", index)
			ms.Equal("First or last names should not be empty", err.Error(), "Case#: ", index)

			continue
		}

		ms.NoError(err, "Case#: %v", index)
		ms.Equal(tc.ExpectedResult, fullName, "Case#: %v", index)
	}
}
