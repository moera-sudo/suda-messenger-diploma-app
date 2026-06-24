// GENERATED CODE - DO NOT MODIFY BY HAND
// dart format width=80

// **************************************************************************
// InjectableConfigGenerator
// **************************************************************************

// ignore_for_file: type=lint
// coverage:ignore-file

// ignore_for_file: no_leading_underscores_for_library_prefixes
import 'package:dio/dio.dart' as _i361;
import 'package:flutter_secure_storage/flutter_secure_storage.dart' as _i558;
import 'package:get_it/get_it.dart' as _i174;
import 'package:injectable/injectable.dart' as _i526;
import 'package:shared_preferences/shared_preferences.dart' as _i460;

import '../../features/auth/data/repositories/auth_repository_impl.dart'
    as _i153;
import '../../features/auth/domain/repositories/auth_repository.dart' as _i787;
import '../../features/auth/presentation/bloc/auth_bloc.dart' as _i797;
import '../../features/channels/data/repositories/channel_repository_impl.dart'
    as _i585;
import '../../features/channels/domain/repositories/channel_repository.dart'
    as _i236;
import '../../features/channels/presentation/bloc/treasury_cubit.dart' as _i848;
import '../../features/chat/data/repositories/chat_repository_impl.dart'
    as _i504;
import '../../features/chat/domain/repositories/chat_repository.dart' as _i420;
import '../../features/chat/presentation/bloc/chat_detail_bloc.dart' as _i57;
import '../../features/chat/presentation/bloc/chat_list_bloc.dart' as _i2;
import '../../features/chat/presentation/bloc/socket_bloc.dart' as _i16;
import '../../features/friends/data/repositories/friends_repository_impl.dart'
    as _i120;
import '../../features/friends/domain/repositories/friends_repository.dart'
    as _i30;
import '../../features/friends/presentation/bloc/friends_cubit.dart' as _i877;
import '../../features/media/data/repositories/media_repository_impl.dart'
    as _i110;
import '../../features/media/domain/repositories/media_repository.dart'
    as _i459;
import '../../features/preferences/data/repositories/preferences_repository_impl.dart'
    as _i450;
import '../../features/preferences/domain/repositories/preferences_repository.dart'
    as _i329;
import '../../features/preferences/presentation/bloc/preferences_cubit.dart'
    as _i460;
import '../../features/profile/data/repositories/profile_repository_impl.dart'
    as _i334;
import '../../features/profile/domain/repositories/profile_repository.dart'
    as _i894;
import '../../features/profile/presentation/bloc/profile_bloc.dart' as _i469;
import '../../features/search/data/repositories/search_repository_impl.dart'
    as _i1017;
import '../../features/search/domain/repositories/search_repository.dart'
    as _i357;
import '../../features/search/presentation/bloc/chat_search_bloc.dart' as _i608;
import '../../features/search/presentation/bloc/search_bloc.dart' as _i552;
import '../../features/user/data/repositories/user_repository_impl.dart'
    as _i664;
import '../../features/user/domain/repositories/user_repository.dart' as _i237;
import '../../features/user/presentation/bloc/contacts_bloc.dart' as _i436;
import '../../features/user/presentation/bloc/user_profile_bloc.dart' as _i322;
import '../../features/wallet/data/repositories/wallet_repository_impl.dart'
    as _i690;
import '../../features/wallet/domain/repositories/wallet_repository.dart'
    as _i571;
import '../../shared/data/api/api_client.dart' as _i513;
import '../../shared/data/api/socket_client.dart' as _i997;
import '../../shared/data/auth/auth_event_bus.dart' as _i419;
import '../../shared/data/auth/current_user.dart' as _i892;
import '../../shared/data/contacts/device_contacts_service.dart' as _i919;
import '../../shared/data/dio/data/auth_interceptor.dart' as _i976;
import '../../shared/data/dio/data/dio_module.dart' as _i301;
import '../../shared/data/locale/shared_storage.dart' as _i684;
import '../../shared/data/logger/app_logger.dart' as _i348;
import '../../shared/data/storage/hive_storage.dart' as _i798;
import '../../shared/data/storage/secure_storage_client.dart' as _i130;
import '../../shared/presentation/bloc/settings_cubit.dart' as _i568;
import '../config/app_config.dart' as _i650;
import '../navigation/app_router.dart' as _i630;
import 'register_module.dart' as _i291;

extension GetItInjectableX on _i174.GetIt {
  // initializes the registration of main-scope dependencies inside of GetIt
  Future<_i174.GetIt> init({
    String? environment,
    _i526.EnvironmentFilter? environmentFilter,
  }) async {
    final gh = _i526.GetItHelper(this, environment, environmentFilter);
    final registerModule = _$RegisterModule();
    final dioModule = _$DioModule();
    await gh.factoryAsync<_i460.SharedPreferences>(
      () => registerModule.prefs,
      preResolve: true,
    );
    gh.singleton<_i650.AppConfig>(() => _i650.AppConfig());
    gh.lazySingleton<_i558.FlutterSecureStorage>(
      () => registerModule.secureStorage,
    );
    gh.lazySingleton<_i630.AppRouter>(() => _i630.AppRouter());
    gh.lazySingleton<_i419.AuthEventBus>(
      () => _i419.AuthEventBus(),
      dispose: (i) => i.dispose(),
    );
    gh.lazySingleton<_i919.DeviceContactsService>(
      () => _i919.DeviceContactsService(),
    );
    gh.lazySingleton<_i348.AppLogger>(() => _i348.AppLogger());
    gh.lazySingleton<_i798.HiveClient>(
      () => _i798.HiveClient(gh<_i348.AppLogger>()),
    );
    gh.lazySingleton<_i684.PreferencesLocalSource>(
      () => _i684.PreferencesLocalSource(gh<_i460.SharedPreferences>()),
    );
    gh.lazySingleton<_i130.SecureStorageClient>(
      () => _i130.SecureStorageClient(
        gh<_i558.FlutterSecureStorage>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i892.CurrentUser>(
      () => _i892.CurrentUser(gh<_i130.SecureStorageClient>()),
    );
    gh.lazySingleton<_i997.SocketClient>(
      () => _i997.SocketClient(
        gh<_i130.SecureStorageClient>(),
        gh<_i348.AppLogger>(),
        gh<_i650.AppConfig>(),
      ),
    );
    gh.factory<_i16.SocketBloc>(
      () => _i16.SocketBloc(gh<_i997.SocketClient>()),
    );
    gh.factory<_i568.SettingsCubit>(
      () => _i568.SettingsCubit(gh<_i684.PreferencesLocalSource>()),
    );
    gh.factory<_i976.AuthInterceptor>(
      () => _i976.AuthInterceptor(
        gh<_i130.SecureStorageClient>(),
        gh<_i348.AppLogger>(),
        gh<_i419.AuthEventBus>(),
        gh<_i650.AppConfig>(),
      ),
    );
    gh.lazySingleton<_i361.Dio>(
      () => dioModule.dio(gh<_i650.AppConfig>(), gh<_i976.AuthInterceptor>()),
    );
    gh.lazySingleton<_i513.ApiClient>(() => _i513.ApiClient(gh<_i361.Dio>()));
    gh.lazySingleton<_i357.SearchRepository>(
      () => _i1017.SearchRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i329.PreferencesRepository>(
      () => _i450.PreferencesRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i787.AuthRepository>(
      () => _i153.AuthRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i130.SecureStorageClient>(),
      ),
    );
    gh.lazySingleton<_i237.UserRepository>(
      () => _i664.UserRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i797.AuthBloc>(
      () =>
          _i797.AuthBloc(gh<_i787.AuthRepository>(), gh<_i419.AuthEventBus>()),
    );
    gh.lazySingleton<_i459.MediaRepository>(
      () => _i110.MediaRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i30.FriendsRepository>(
      () => _i120.FriendsRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i571.WalletRepository>(
      () => _i690.WalletRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i608.ChatSearchBloc>(
      () => _i608.ChatSearchBloc(
        gh<_i357.SearchRepository>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i236.ChannelRepository>(
      () => _i585.ChannelRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i420.ChatRepository>(
      () => _i504.ChatRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i877.FriendsCubit>(
      () => _i877.FriendsCubit(
        gh<_i30.FriendsRepository>(),
        gh<_i997.SocketClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i436.ContactsBloc>(
      () =>
          _i436.ContactsBloc(gh<_i237.UserRepository>(), gh<_i348.AppLogger>()),
    );
    gh.factory<_i322.UserProfileBloc>(
      () => _i322.UserProfileBloc(
        gh<_i237.UserRepository>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.lazySingleton<_i894.ProfileRepository>(
      () => _i334.ProfileRepositoryImpl(
        gh<_i513.ApiClient>(),
        gh<_i130.SecureStorageClient>(),
        gh<_i459.MediaRepository>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i460.PreferencesCubit>(
      () => _i460.PreferencesCubit(
        gh<_i329.PreferencesRepository>(),
        gh<_i997.SocketClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i469.ProfileBloc>(
      () => _i469.ProfileBloc(
        gh<_i894.ProfileRepository>(),
        gh<_i571.WalletRepository>(),
      ),
    );
    gh.factory<_i552.SearchBloc>(
      () => _i552.SearchBloc(
        gh<_i357.SearchRepository>(),
        gh<_i420.ChatRepository>(),
      ),
    );
    gh.factory<_i2.ChatListBloc>(
      () => _i2.ChatListBloc(
        gh<_i420.ChatRepository>(),
        gh<_i997.SocketClient>(),
        gh<_i130.SecureStorageClient>(),
        gh<_i348.AppLogger>(),
      ),
    );
    gh.factory<_i57.ChatDetailBloc>(
      () => _i57.ChatDetailBloc(
        gh<_i420.ChatRepository>(),
        gh<_i459.MediaRepository>(),
        gh<_i997.SocketClient>(),
        gh<_i130.SecureStorageClient>(),
        gh<_i348.AppLogger>(),
        gh<_i329.PreferencesRepository>(),
      ),
    );
    gh.factory<_i848.TreasuryCubit>(
      () => _i848.TreasuryCubit(gh<_i236.ChannelRepository>()),
    );
    return this;
  }
}

class _$RegisterModule extends _i291.RegisterModule {}

class _$DioModule extends _i301.DioModule {}
