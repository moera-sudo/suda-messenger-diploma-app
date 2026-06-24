import 'package:equatable/equatable.dart';

class MyProfile extends Equatable {
  final String id;
  final String username;
  final String email;
  final String displayName;
  final String firstName;
  final String lastName;
  final String bio;
  final String? avatarUrl;
  final bool isOnline;

  const MyProfile({
    required this.id,
    required this.username,
    required this.email,
    required this.displayName,
    required this.firstName,
    required this.lastName,
    required this.bio,
    this.avatarUrl,
    this.isOnline = false,
  });

  MyProfile copyWith({
    String? id,
    String? username,
    String? email,
    String? displayName,
    String? firstName,
    String? lastName,
    String? bio,
    String? avatarUrl,
    bool? isOnline,
  }) {
    return MyProfile(
      id: id ?? this.id,
      username: username ?? this.username,
      email: email ?? this.email,
      displayName: displayName ?? this.displayName,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      bio: bio ?? this.bio,
      avatarUrl: avatarUrl ?? this.avatarUrl,
      isOnline: isOnline ?? this.isOnline,
    );
  }

  @override
  List<Object?> get props => [
    id,
    username,
    email,
    displayName,
    firstName,
    lastName,
    bio,
    avatarUrl,
    isOnline,
  ];
}

class UpdateMyProfilePayload {
  final String displayName;
  final String firstName;
  final String lastName;
  final String bio;

  const UpdateMyProfilePayload({
    required this.displayName,
    required this.firstName,
    required this.lastName,
    required this.bio,
  });

  Map<String, dynamic> toJson() {
    return {
      'display_name': displayName,
      'first_name': firstName,
      'last_name': lastName,
      'bio': bio,
    };
  }
}
