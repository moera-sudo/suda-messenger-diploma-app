import 'package:hive_flutter/hive_flutter.dart';
import 'package:injectable/injectable.dart';
import '../../data/logger/app_logger.dart';


@lazySingleton
class HiveClient {
  final AppLogger _logger;
  HiveClient(this._logger);

  Future<Box<T>> _openBox<T>(String boxName) async {
    if (Hive.isBoxOpen(boxName)) {
      return Hive.box<T>(boxName);
    } else {
      _logger.trace("Opening Hive Box: $boxName");
      return await Hive.openBox<T>(boxName);
    }
  }

  Future<void> put<T>(String boxName, String key, T value) async {
    try {
      final box = await _openBox<T>(boxName);
      await box.put(key, value);
      _logger.trace("Hive Put [$boxName] key: $key");
    } catch(e) {
      _logger.error("Hive put errror in $boxName", e);
      throw Exception("Local storage write error");
    }
  } 

  Future<T?> get<T>(String boxName, String key, {T? defaultValue}) async{
    try{
      final box = await _openBox<T>(boxName);
      return box.get(key, defaultValue: defaultValue);
    } catch(e) {
      _logger.error("Hive Get error in $boxName", e);
      return defaultValue;
    }
  }

  Future<List<T>> getAll<T>(String boxName) async{
    try{
      final box = await _openBox<T>(boxName);
      return box.values.toList();
    } catch(e) {
      _logger.error("Hive get all error in $boxName", e);
      return [];
    }
  }

  Future<void> delete<T>(String boxName, String key) async {
    try {
      final box = await _openBox<T>(boxName);
      await box.delete(key);
      _logger.trace("Hive Delete [$boxName] key: $key");
    } catch (e) {
      _logger.error("Hive Delete Error in $boxName", e);
    }
  }

  Future<void> clear<T>(String boxName) async {
    try {
      final box = await _openBox<T>(boxName);
      await box.clear();
      _logger.warning("Hive Box Cleared: $boxName");
    } catch (e) {
      _logger.error("Hive Clear Error in $boxName", e);
    }
  }

  Future<void> deleteBoxFromDisk(String boxName) async {
     try {
       // Сначала убедимся, что она закрыта или открыта, API требует разного подхода
       // Но проще просто вызвать deleteFromDisk если она открыта
       if (Hive.isBoxOpen(boxName)) {
         await Hive.box(boxName).deleteFromDisk();
       } else {
         await Hive.deleteBoxFromDisk(boxName);
       }
       _logger.warning("Hive Box Deleted From Disk: $boxName");
     } catch (e) {
       _logger.error("Hive DeleteBox Error", e);
     }
  }

}

