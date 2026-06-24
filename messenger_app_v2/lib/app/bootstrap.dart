import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter/foundation.dart';
import 'package:hive_flutter/adapters.dart';
import 'package:messenger_app_v2/app/DI/get_it.dart';

import 'lifespan.dart';
import 'app.dart';
import '../firebase_options.dart';
import '../shared/data/logger/app_logger.dart';
import '../shared/data/notifications/notification_service.dart';
import '../shared/data/storage/secure_storage_client.dart';
import '../shared/data/api/api_client.dart';

Future<void> bootstrap() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
  FirebaseMessaging.onBackgroundMessage(firebaseBackgroundMessageHandler);

  await Hive.initFlutter();

  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
  ]);

  // Инициализация DI
  await configureDependencies();

  final logger = sl<AppLogger>();
  logger.info("Application starting...");

  await NotificationService.setup();

  if (kDebugMode) {
    final storage = sl<SecureStorageClient>();
    await storage.clearStorageForDebug();
  }

  try {
    final api = sl<ApiClient>();
    logger.info("Trying to send a request");

    await api.get("/health");

    logger.info("Online D:");
  } catch (e) {
    logger.fatal("Server is offline");
  }

  WidgetsBinding.instance.addObserver(Lifespan(logger));

  runApp(const MainApp());
}